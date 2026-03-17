<template>
  <div class="relative w-full h-4 flex items-center justify-center px-2">
    <div class="flex-1 flex justify-start pl-2">
      <span v-if="sum > 0" class="text-xs text-blue-500 font-semibold bg-blue-50 px-2 py-1 rounded min-w-3.5rem text-center tabular-nums">
        {{ index }} / {{ sum }}
      </span>
    </div>

    <div class="flex-3 flex justify-center w-full min-w-0 px-4">
      <div class="flex items-center text-gray-700 font-medium max-w-full">
        <!-- 无后缀的文件主体，加 truncate 以显示省略号 -->
        <span class="truncate block">
          {{ currentName }}
        </span>
        <!-- 文件后缀，加 shrink-0 防止其被挤压截断 -->
        <span class="shrink-0">
          {{ currentExtension }}
        </span>
      </div>
    </div>

    <div class="flex-1 flex justify-end gap-1 pr-2">
      <button 
        @click="prevTrack" 
        :disabled="sum === 0 || index === 1"
        class="flex items-center justify-center w-6 h-6 rounded-md bg-gray-100 hover:bg-blue-100 text-gray-500 hover:text-blue-600 disabled:opacity-30 disabled:cursor-not-allowed transition-colors"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z" clip-rule="evenodd" />
        </svg>
      </button>

      <button 
        @click="nextTrack" 
        :disabled="sum === 0 || index === sum"
        class="flex items-center justify-center w-6 h-6 rounded-md bg-gray-100 hover:bg-blue-100 text-gray-500 hover:text-blue-600 disabled:opacity-30 disabled:cursor-not-allowed transition-colors"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd" />
        </svg>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { tracks, currentIndex } from '../../store/audioManager';

const currentTrack = computed(() => {
    const track = tracks.value[currentIndex.value];
    return track ? track.filename : null;
})

const currentName = computed(() => {
    const track = tracks.value[currentIndex.value];
    if (!track || !track.filename) return 'No track loaded';
    
    const lastDotIndex = track.filename.lastIndexOf('.');
    if (lastDotIndex <= 0) return track.filename;
    
    return track.filename.substring(0, lastDotIndex);
})

const currentExtension = computed(() => {
    const track = tracks.value[currentIndex.value];
    if (!track || !track.filename) return '';
    
    const lastDotIndex = track.filename.lastIndexOf('.');
    if (lastDotIndex <= 0) return '';
    
    return track.filename.substring(lastDotIndex); // 包含点
})

const index = computed(() => {
    return tracks.value.length > 0 ? currentIndex.value + 1 : 0;
})

const sum = computed(() => {
    return tracks.value.length;
})


const prevTrack = () => {
    if (sum.value > 0 && currentIndex.value > 0) {
        currentIndex.value--;
    }
}

const nextTrack = () => {
    if (sum.value > 0 && currentIndex.value < sum.value - 1) {
        currentIndex.value++;
    }
}
</script>