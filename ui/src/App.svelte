<script lang="ts">
  import "./app.css";
  import { showPopup, togglePopup } from "./store";
  import axios from "axios";
  import { Button } from "$lib/components/ui/button";
  import UserPost from "$lib/components/ui/userpost/userpost.svelte";
  import UploadPost from "$lib/components/ui/uploadpost/uploadpost.svelte";

  interface Likes {
    id: number;
    user_id: number;
    post_id: number;
    created_at: Date;
  }

  interface Comments {
    id: number;
    user_id: number;
    post_id: number;
    content: string;
    created_at: Date;
  }


  interface Post {
    id: number;
    user_id: number;
    media: Base64URLString;
    caption: string;
    upload_time: Date;
    likes: Likes[];
    comments: Comments[];
  }

  let posts: Post[] = [];

  const fetchPosts = async () => {
  try {
    const response = await axios.get<{ Posts: Post[] }>("http://localhost:8080/posts");
    posts = response.data.Posts.map(post => ({
      ...post,
      media: post.media.startsWith("data:image/")
        ? post.media
        : "data:image/png;base64," + post.media,
    }));
    console.log("Posts fetched:", posts);
  } catch (error) {
    console.error("Error fetching posts:", error);
  }
};

function checkTime(upload_time: any): string {
  if (typeof upload_time === 'string') {
    upload_time = new Date(upload_time);
  } else if (typeof upload_time === 'number') {
    upload_time = new Date(upload_time * 1000);
  }

  let diff = Date.now() - upload_time.getTime();

  const minute = 60 * 1000;
  const hour = 60 * minute;
  const day = 24 * hour;

  if (diff < minute) {
    return "Just now";
  } else if (diff < hour) {
    let minutes = Math.floor(diff / minute);
    return `${minutes} minute${minutes > 1 ? 's' : ''} ago`;
  } else if (diff < day) {
    let hours = Math.floor(diff / hour);
    return `${hours} hour${hours > 1 ? 's' : ''} ago`;
  } else {
    let days = Math.floor(diff / day);
    return `${days} day${days > 1 ? 's' : ''} ago`;
  }
}


</script>

<Button on:click={togglePopup}>Upload</Button>
{#if $showPopup}
  <UploadPost />
{/if}
<button on:click={fetchPosts}>Show posts</button>

{#each posts as post}
<UserPost
  user_id={post.user_id}
  media={post.media}
  caption={post.caption}
  upload_time={checkTime(post.upload_time)}
  likes={post.likes}
  comments={post.comments}
/>
{/each}


