<script lang="ts" setup>
import { ref , onMounted, onUnmounted} from 'vue'
import { OnFileDrop , OnFileDropOff} from '../wailsjs/runtime/runtime'

import { GetFileInfo } from '../wailsjs/go/main/App'

interface AudioTrack {
  name: string;
  path: string;
  size: number;

  edits: {
    volume: number;
    trimStart: number;
    trimEnd: number;
  }
}

const isDragging = ref(false)
const fileList = ref<Array<{ name: string; size: number; path: string } | null>>([])


onUnmounted(() => {
  OnFileDropOff()
})

</script>

<template>
  

</template>

<style scoped>

</style>
