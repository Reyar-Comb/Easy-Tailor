<template>
    <div class="w-full h-32 bg-gray-200 flex items-center justify-center relative">
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
                    :key="`${tick.x}-${tick.label}`"
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
}) 

const formatTime = (seconds: number) => {
    if (seconds < 60) {
        return `0:${Math.floor(seconds).toString().padStart(2, '0')}`;
    }

    const mins = Math.floor(seconds / 60);
    const secs = Math.floor(seconds % 60);
    return `${mins}:${secs.toString().padStart(2, '0')}`;
};

const ticks = computed(() => {
    const durationSec = (props.durationMs ?? 0) / 1000;
    if (!durationSec || durationSec <= 0) return [];

    const targetTickCount = 6;
    const rawStep = durationSec / targetTickCount;
    const candidates = [1, 2, 5, 10, 15, 30, 60, 120, 300, 600, 900, 1800];
    const step = candidates.find(candidate => candidate >= rawStep) ?? candidates[candidates.length - 1];

    const result: { x: number; label: string }[] = [];
    for (let t = 0; t <= durationSec; t += step) {
        result.push({
            x: Math.min(100, (t / durationSec) * 100),
            label: formatTime(t),
        });
    }

    if (result[result.length - 1]?.x !== 100) {
        result.push({
            x: 100,
            label: formatTime(durationSec),
        });
    }

    return result;
});

const canvasRef = ref<HTMLCanvasElement | null>(null);
const drawWaveform = () => {
    const canvas = canvasRef.value;
    if (!canvas || !props.data) return;

    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    const rect = canvas.getBoundingClientRect();
    const dpr = window.devicePixelRatio || 1;

    canvas.width = rect.width * dpr;
    canvas.height = rect.height * dpr;

    ctx.scale(dpr, dpr);
    ctx.clearRect(0, 0, rect.width, rect.height);

    ctx.beginPath();
    ctx.moveTo(0, rect.height / 2);
    ctx.lineTo(rect.width, rect.height / 2);
    ctx.strokeStyle = '#bfdbfe';
    ctx.lineWidth = 1;
    ctx.stroke();

    ctx.fillStyle = '#3b82f6';

    const dataLength = props.data!.length;
    const barWidth = rect.width / dataLength;

    const maxBarHeight = (rect.height / 2) * 0.9;

    for (let i = 0; i < dataLength; i++) {
        const x = i * barWidth;
        
        const barW = Math.max(1, barWidth * 0.8);

        const ampL = props.data![i][0];
        const heightL = ampL * maxBarHeight;
        ctx.fillRect(x, rect.height / 2, barW, heightL);

        const ampR = props.data![i][1];
        const heightR = ampR * maxBarHeight;
        ctx.fillRect(x, rect.height / 2 - heightR, barW, heightR);
    }
};

watch(() => props.data, () => {
    nextTick(() => {
        drawWaveform();
    });
}, { deep: true });

onMounted(() => {
    window.addEventListener('resize', drawWaveform);
    if (hasData.value) {
        drawWaveform();
    }
});

onBeforeUnmount(() => {
    window.removeEventListener('resize', drawWaveform);
});
</script>
