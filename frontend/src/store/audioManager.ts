import { ref, reactive } from 'vue'

export interface AudioTrack {
  id: string;          
  filename: string;    
  filepath: string;         
  waveform: number[];
  edits:
  {
    volume: number;
    pitch: number;
    TrimStart: number;
    TrimEnd: number;
  }
}

export const tracks = ref<AudioTrack[]>([]);
export const currentIndex = ref<number>(0);

export const currentTrack = () => tracks.value[currentIndex.value] || null;
export const addTrack = (track: AudioTrack) => {
    tracks.value.push(track);
    currentIndex.value = tracks.value.length - 1;
}
export const cancelTrack = (index: number) => {
    if (index < 0 || index >= tracks.value.length) return;
    tracks.value.splice(index, 1);
    if (currentIndex.value >= tracks.value.length) {
        currentIndex.value = tracks.value.length - 1;
    }
}