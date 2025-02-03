<script lang="ts">
  import { showPopup, togglePopup, notifications, fetchNotifications, markNotificationsAsSeen } from "$lib/../store";
  import axios from "axios";
  import { Button } from "$lib/components/ui/button";
  import UserPost from "$lib/components/ui/userpost/userpost.svelte";
  import UploadPost from "$lib/components/ui/uploadpost/uploadpost.svelte";
  import LogoutButton from "$lib/components/ui/logoutButton/logoutButton.svelte";
  import { onMount } from "svelte";
  import { currentUsername, currentUserId, fetchCurrentUser } from "./userstore";
  import { get } from "svelte/store";
  import { Link } from 'svelte-routing';

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
  let usernameCache = new Map<number, string>();

  // Fetch posts
  const fetchPosts = async () => {
    try {
      const response = await axios.get<{ Posts: Post[] }>("http://4.234.181.167:8080/posts");
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

  // Fetch username by user_id and store in cache
  const fetchUsername = async (user_id: number): Promise<string> => {
    if (usernameCache.has(user_id)) {
      return usernameCache.get(user_id)!; // Return cached username
    }

    try {
      const response = await axios.get(`http://4.234.181.167:8080/users/username/${user_id}`);
      const username = response.data.data;
      usernameCache.set(user_id, username); // Cache the username
      return username;
    } catch (error) {
      console.error("Error fetching username:", error);
      return "Unknown User"; // Fallback username
    }
  };

  // Format upload time
  function checkTime(upload_time: any): string {
    if (typeof upload_time === "string") {
      upload_time = new Date(upload_time);
    } else if (typeof upload_time === "number") {
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
      return `${minutes} minute${minutes > 1 ? "s" : ""} ago`;
    } else if (diff < day) {
      let hours = Math.floor(diff / hour);
      return `${hours} hour${hours > 1 ? "s" : ""} ago`;
    } else {
      let days = Math.floor(diff / day);
      return `${days} day${days > 1 ? "s" : ""} ago`;
    }
  }

  onMount(async () => {
    await fetchCurrentUser();
    await fetchPosts();
    await fetchNotifications();
    await markNotificationsAsSeen();
  });

  function handlePostUploaded() {
    fetchPosts();
  }
</script>

<style>
  #notifications {
    position: fixed;
    top: 10px;
    right: 10px;
    width: 300px;
    max-height: 400px;
    overflow-y: auto;
    background: white;
    border: 1px solid #ddd;
    border-radius: 8px;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    padding: 10px;
  }

  #notifications h3 {
    margin-top: 0;
  }

  #notifications ul {
    list-style: none;
    padding: 0;
  }

  #notifications li {
    margin-bottom: 10px;
    padding: 10px;
    border-bottom: 1px solid #ddd;
  }

  #notifications li:last-child {
    border-bottom: none;
  }

  #fixedButtons {
    position: fixed;
    top: 20px;
    left: 20px;
    display: flex;
    flex-direction: column;
    gap: 10px;
    align-items: center; 
  }

  #userPosts {
    margin-top:10px;
  }
  #username {
    font-weight: bold;
  }

</style>

<div id="fixedButtons">
  <Link to="/" class="mb-4 text-[20px] font-bold cursor-pointer">Fakestagram</Link>
  <Button 
  class="flex items-center justify-center px-5 py-2.5 bg-gray-100 text-gray-800 rounded-md shadow-md transition duration-300 ease-in-out hover:bg-gray-200 hover:text-black"
  on:click={togglePopup}
>
  <span class="text-[16px] font-bold">Upload</span>
</Button>
</div>

{#if $showPopup}
  <UploadPost on:postUploaded={handlePostUploaded}/>
{/if}

<div id="notifications">
  <h3>Notifications  <Link to={`/${get(currentUsername)}`}><span id="username">{get(currentUsername)}</span></Link></h3>
  <ul>
    {#each $notifications as notification}
      <li>
        <p>{notification.content}</p>
        <small>{new Date(notification.created_at).toLocaleString()}</small>
      </li>
    {/each}
  </ul>
</div>

<LogoutButton />

<div id="userPosts">
{#each posts as post}
  {#await fetchUsername(post.user_id) then username}
    <UserPost
    user_id={post.user_id}
    user_name={username}
    media={post.media}
    caption={post.caption}
    upload_time={checkTime(post.upload_time)}
    post_id={post.id}
    CurrentUserId={get(currentUserId)}
    CurrentUsername={get(currentUsername)}
    />
  {/await}
{/each}
</div>