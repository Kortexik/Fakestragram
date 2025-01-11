import { writable } from "svelte/store";

export const showPopup = writable(false);

export const togglePopup = () => {
    showPopup.update((visible) => !visible);
  };