<template>
  <button 
    class="bg-blue-500 hover:bg-blue-600 text-white font-medium py-2 px-4 rounded shadow-sm transition-colors w-20 flex justify-center items-center"
    @click="handleClick"
  >
    <slot></slot> 
  </button>
</template>

<script setup lang="ts">

const handleClick = () => {
  console.log("按钮被点击了")
}
</script>