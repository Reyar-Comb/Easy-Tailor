<template>
    <div class="w-full h-32 bg-gray-200 flex items-center justify-center">
        <div v-if="!hasData" class="absolute inset-0 flex items-center justify-center text-blue-300 font-medium select-none">
            No waveform data available
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
}>();

const hasData = computed(() => {
    return props.data && props.data.length > 0;
}) 

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

    ctx.fillStyle = '#1d4ed8';

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