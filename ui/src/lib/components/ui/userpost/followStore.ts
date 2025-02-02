import { writable } from "svelte/store";

export const followState = writable<{ [key: number]: boolean }>({});

export const setFollowState = (userId: number, isFollowing: boolean) => {
  followState.update((state) => {
    state[userId] = isFollowing;
    return state;
  });
};