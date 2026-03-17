<script lang="ts" setup>
import { ref , onMounted, onUnmounted} from 'vue';
import { OnFileDrop , OnFileDropOff} from '../wailsjs/runtime/runtime';

import { addTrack , currentIndex , tracks } from './store/audioManager';

import { GetFileInfo } from '../wailsjs/go/main/App'
import ActionButton from './components/TopPanel/ActionButton.vue';
import FileNavigator from './components/TopPanel/FileNavigator.vue';

const IsDragging = ref(false);

// Handle Percentage
const topHeightPercentage = ref(60); // 默认上方占 60%
const isResizing = ref(false);

const startResize = (e: MouseEvent) => {
  isResizing.value = true;
  document.body.style.cursor = 'row-resize'; 
};

const doResize = (e: MouseEvent) => {
  if (!isResizing.value) return;
  let newPercentage = (e.clientY / window.innerHeight) * 100;
  
  if (newPercentage < 20) newPercentage = 20;
  if (newPercentage > 80) newPercentage = 80;

  topHeightPercentage.value = newPercentage;
};

const stopResize = () => {
  isResizing.value = false;
  document.body.style.cursor = ''; 
};


//Handle Mounted
onMounted(() => {
  window.addEventListener('mousemove', doResize);
  window.addEventListener('mouseup', stopResize);


  OnFileDrop(async (x: number, y: number, paths: string[]) => {
    // Get All File Info by order
    IsDragging.value = false;
    for (const filepath of paths) {
      const info = await GetFileInfo(filepath)
      if (info) {
        addTrack({
          id: Date.now().toString(),          
          filename: info.name, 
          filepath: info.path,      
          waveform: [],
          edits:
          {
            volume: 50,
            pitch: 50,
            TrimStart: 0,
            TrimEnd: -1
          }
        })
      }
    }
    console.log('文件已添加:', paths)
  }, true);
});

onUnmounted(() => {
  window.removeEventListener('mousemove', doResize);
  window.removeEventListener('mouseup', stopResize);

  OnFileDropOff()
})

</script>

<template>
  <div class="w-screen h-screen flex flex-col bg-gray-50 text-slate-800"
    :class="{ 'select-none': isResizing }"
    style="--wails-drop-target: drop;"
    @dragover.prevent="IsDragging = true"
    @dragleave="IsDragging = false"
    @drop.prevent="IsDragging = false">
    
    <div 
      class="relative border-b-0 flex flex-row items-stretch p-3 gap-4"
      :style="{ height: topHeightPercentage + '%' }"
    >    
        <div class="flex flex-col justify-center items-center gap-3 w-28 shrink-0">
          <ActionButton>Import</ActionButton>
          <ActionButton>Settings</ActionButton>
        </div>

        <div class="flex flex-col flex-1 min-w-0">
          <div class="mb-1">
            <FileNavigator />
          </div>

          <div class="flex-1 w-full relative bg-blue-50 border border-blue-200 rounded-lg shadow-inner overflow-hidden">
            <WaveformView />
          </div>
        </div>

        <div class="flex flex-col justify-center items-center gap-3 w-28 shrink-0">
          <ActionButton>Export</ActionButton>
          <ActionButton>Cancel</ActionButton>
        </div>
        
    </div>

    <div 
      class="w-full flex justify-center items-center h-2 -my-1 z-10 cursor-row-resize group"
      @mousedown.prevent="startResize"
    >
      <div class="w-full h-2px bg-blue-100 group-hover:bg-blue-400 group-hover:h-1 transition-all duration-150"></div>
    </div>


    <div class="flex-1 bg-white p-6 grid grid-cols-4 gap-4 overflow-y-auto">
    </div>
  </div>

  <div 
    v-if="IsDragging" 
    class="fixed inset-0 z-100 border-6px bg-blue-500/40 backdrop-blur-[2px] flex items-center justify-center border-4 border-blue-400 pointer-events-none"
  >
    <span class="text-white/70 text-3xl tracking-widest">
      DROP AUDIOS HERE
    </span>
  </div>
</template>

<style scoped>

</style>
