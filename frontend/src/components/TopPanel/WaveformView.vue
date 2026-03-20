<template>
    <div
        ref="containerRef"
        class="w-full h-32 bg-gray-200 flex items-center justify-center relative overflow-hidden"
        :class="{
            'cursor-grab': hasData,
            'cursor-grabbing': isPanning,
        }"
        @wheel.prevent="handleWheel"
        @pointerdown="handlePointerDown"
        @pointermove="handlePointerMove"
        @pointerup="handlePointerUp"
        @pointercancel="handlePointerUp"
        @pointerleave="handlePointerUp"
    >
        <div v-if="!hasData" class="absolute inset-0 flex items-center justify-center text-blue-300 font-medium select-none">
            No waveform data available
        </div>

        <div
            v-if="hasData && ticks.length > 0"
            class="absolute left-0 right-0 top-0 h-3 px-1 pointer-events-none z-10"
        >
            <div class="relative w-full h-full border-b border-blue-200/80 bg-white/45 backdrop-blur-[1px]">
                <div
                    v-for="tick in ticks"
                    :key="`${tick.index}-${tick.label}`"
                    class="absolute top-0 bottom-0"
                    :style="{ left: `${tick.x}%` }"
                >
                    <div class="w-px h-0.5 bg-blue-300"></div>
                    <div class="text-[8px] leading-none text-blue-700 -translate-x-1/2 mt-0 tabular-nums whitespace-nowrap">
                        {{ tick.label }}
                    </div>
                </div>
            </div>
        </div>

        <canvas
            ref="canvasRef"
            class="w-full h-full absolute inset-0 block"
            :class="{ 'opacity-0': !hasData }"
        ></canvas>
    </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onBeforeUnmount, nextTick, computed } from 'vue';

const props = defineProps<{
    data?: number[][] | null;
    durationMs?: number;
}>();

const hasData = computed(() => {
    return props.data && props.data.length > 0;
});

const containerRef = ref<HTMLDivElement | null>(null);
const canvasRef = ref<HTMLCanvasElement | null>(null);

const zoom = ref(1);
const minZoom = 1;
const maxZoom = 32;
const zoomStep = 1.2;

const viewStartMs = ref(0);

const isPanning = ref(false);
const panStartClientX = ref(0);
const panStartViewStartMs = ref(0);

const totalDurationMs = computed(() => props.durationMs ?? 0);
const viewDurationMs = computed(() => {
    if (totalDurationMs.value <= 0) return 0;
    return totalDurationMs.value / zoom.value;
});

const maxViewStartMs = computed(() => {
    return Math.max(0, totalDurationMs.value - viewDurationMs.value);
});

const viewEndMs = computed(() => {
    return Math.min(totalDurationMs.value, viewStartMs.value + viewDurationMs.value);
});

const visibleSampleRange = computed(() => {
    if (!props.data?.length || totalDurationMs.value <= 0) {
        return { startIndex: 0, endIndex: 0, count: 0 };
    }

    const totalSamples = props.data.length;
    const msPerSample = totalDurationMs.value / totalSamples;

    const startIndex = Math.max(0, Math.floor(viewStartMs.value / msPerSample));
    const endIndex = Math.min(totalSamples, Math.ceil(viewEndMs.value / msPerSample));
    const count = Math.max(1, endIndex - startIndex);

    return { startIndex, endIndex, count };
});

const clampViewStart = (nextStartMs: number) => {
    return Math.min(Math.max(0, nextStartMs), maxViewStartMs.value);
};

const resetViewport = () => {
    zoom.value = 1;
    viewStartMs.value = 0;
};

const updateZoomAt = (nextZoom: number, anchorRatio: number) => {
    if (totalDurationMs.value <= 0) return;

    const clampedZoom = Math.min(maxZoom, Math.max(minZoom, nextZoom));
    const safeAnchorRatio = Math.min(Math.max(anchorRatio, 0), 1);
    const anchorMs = viewStartMs.value + safeAnchorRatio * viewDurationMs.value;
    const nextViewDurationMs = totalDurationMs.value / clampedZoom;
    const nextViewStartMs = anchorMs - safeAnchorRatio * nextViewDurationMs;

    zoom.value = clampedZoom;
    viewStartMs.value = clampViewStart(nextViewStartMs);
};

const roundTo = (value: number, digits = 3) => {
    const factor = 10 ** digits;
    return Math.round(value * factor) / factor;
};

const formatTime = (seconds: number, step: number) => {
    const safeSeconds = roundTo(seconds, 3);
    const safeStep = roundTo(step, 3);
    const mins = Math.floor(safeSeconds / 60);
    const secs = Math.floor(safeSeconds % 60);

    if (safeStep >= 1) {
        return `${mins}:${secs.toString().padStart(2, '0')}`;
    }

    if (safeStep >= 0.1) {
        const tenths = Math.floor(roundTo((safeSeconds % 1) * 10, 3));
        return `${mins}:${secs.toString().padStart(2, '0')}.${tenths}`;
    }
    
    const hundredths = Math.floor(roundTo((safeSeconds * 100) % 100, 3));
    return `${mins}:${secs.toString().padStart(2, '0')}.${hundredths.toString().padStart(2, '0')}`;
};

const ticks = computed(() => {
    const visibleDurationSec = viewDurationMs.value / 1000;
    if (!visibleDurationSec || visibleDurationSec <= 0) return [];

    const targetTickCount = 6;
    const rawStep = visibleDurationSec / targetTickCount;
    const candidates = [
        0.05, 0.1, 0.2, 0.5,
        1, 2, 5, 10, 15, 30,
        60, 120, 300, 600, 900, 1800,
    ];
    const step = candidates.find(candidate => candidate >= rawStep) ?? candidates[candidates.length - 1];
    const safeStep = roundTo(step, 3);
    const viewStartSec = roundTo(viewStartMs.value / 1000, 3);
    const viewEndSec = roundTo(viewEndMs.value / 1000, 3);
    const firstTick = roundTo(Math.ceil(viewStartSec / safeStep) * safeStep, 3);

    const result: { index: number; x: number; label: string }[] = [];
    for (let t = firstTick; t <= viewEndSec + safeStep * 0.25; t = roundTo(t + safeStep, 3)) {
        const safeT = roundTo(t, 3);
        result.push({
            index: result.length,
            x: roundTo(Math.min(100, Math.max(0, ((safeT - viewStartSec) / visibleDurationSec) * 100)), 3),
            label: formatTime(safeT, safeStep),
        });
    }
    return result;
});

const drawWaveform = () => {
    const canvas = canvasRef.value;
    if (!canvas) return;

    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    const rect = canvas.getBoundingClientRect();
    const dpr = window.devicePixelRatio || 1;

    canvas.width = rect.width * dpr;
    canvas.height = rect.height * dpr;

    ctx.setTransform(1, 0, 0, 1, 0, 0);
    ctx.scale(dpr, dpr);
    ctx.clearRect(0, 0, rect.width, rect.height);

    if (!props.data?.length) return;

    ctx.beginPath();
    ctx.moveTo(0, rect.height / 2);
    ctx.lineTo(rect.width, rect.height / 2);
    ctx.strokeStyle = '#bfdbfe';
    ctx.lineWidth = 1;
    ctx.stroke();

    ctx.fillStyle = '#3b82f6';

    const maxBarHeight = (rect.height / 2) * 0.9;
    const { startIndex, endIndex, count } = visibleSampleRange.value;
    const barWidth = rect.width / count;

    for (let i = startIndex; i < endIndex; i++) {
        const x = (i - startIndex) * barWidth;
        const barW = Math.max(1, barWidth * 0.8);
        const ampL = props.data[i][0];
        const heightL = ampL * maxBarHeight;
        ctx.fillRect(x, rect.height / 2, barW, heightL);

        const ampR = props.data[i][1];
        const heightR = ampR * maxBarHeight;
        ctx.fillRect(x, rect.height / 2 - heightR, barW, heightR);
    }
};

const handleWheel = (event: WheelEvent) => {
    if (!hasData.value || !containerRef.value) return;

    const rect = containerRef.value.getBoundingClientRect();
    if (rect.width <= 0) return;

    const anchorRatio = (event.clientX - rect.left) / rect.width;
    const direction = event.deltaY < 0 ? zoomStep : 1 / zoomStep;
    updateZoomAt(zoom.value * direction, anchorRatio);
};

const handlePointerDown = (event: PointerEvent) => {
    if (!hasData.value || !containerRef.value) return;

    isPanning.value = true;
    panStartClientX.value = event.clientX;
    panStartViewStartMs.value = viewStartMs.value;
    containerRef.value.setPointerCapture(event.pointerId);
};

const handlePointerMove = (event: PointerEvent) => {
    if (!isPanning.value || !containerRef.value) return;

    const rect = containerRef.value.getBoundingClientRect();
    if (rect.width <= 0) return;

    const deltaX = event.clientX - panStartClientX.value;
    const deltaMs = -(deltaX / rect.width) * viewDurationMs.value;
    viewStartMs.value = clampViewStart(panStartViewStartMs.value + deltaMs);
};

const handlePointerUp = (event: PointerEvent) => {
    if (!isPanning.value || !containerRef.value) return;

    if (containerRef.value.hasPointerCapture(event.pointerId)) {
        containerRef.value.releasePointerCapture(event.pointerId);
    }
    isPanning.value = false;
};

watch(() => [props.data, props.durationMs], () => {
    resetViewport();
    nextTick(drawWaveform);
}, { deep: true });

watch(() => [
    props.durationMs,
    zoom.value,
    viewStartMs.value,
    visibleSampleRange.value.startIndex,
    visibleSampleRange.value.endIndex,
], () => {
    nextTick(drawWaveform);
});

let resizeObserver: ResizeObserver | null = null;

onMounted(() => {
    if (hasData.value) {
        drawWaveform();
    }

    if (containerRef.value) {
        resizeObserver = new ResizeObserver(() => {
            drawWaveform();
        });
        resizeObserver.observe(containerRef.value);
    }
});

onBeforeUnmount(() => {
    resizeObserver?.disconnect();
});
</script>
