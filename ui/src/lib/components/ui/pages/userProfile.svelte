<script lang="ts">
  import { onMount } from "svelte";
  import axios from "axios";
  import { Button } from "$lib/components/ui/button";
  import { currentUsername, currentUserId } from "./userstore";
  import { editing, setEditing } from "$lib/../store"
  import EditProfile from "$lib/components/ui/editProfile/editProfile.svelte";
  import { get } from "svelte/store";

  
  interface User {
    userID: number;
    username: string;
    firstName: string;
    lastName: string;
    bio: string;
    profilePic: Base64URLString;
    createdAt: Date;
    numberOfPosts: number;
    numberOfFollowers: number;
    numberOfFollowees: number;
  };


  let profile: User | null = $state();
  let posts = $state();
  let props = $props();
  let fileSelected = false;
  let isFollowing = $state();
  let wipMessage = $state(false);
  let imageClicked = $state(false);

  function toogleWIPmessage() {
    wipMessage = !wipMessage;
  }

  function toogleImageClicked() {
    imageClicked = !imageClicked;
  }

  function clickFileInput() {
    document.getElementById('file-input').click();
  }
  
function handleFileInputChange(event) {
  selectedFile = event.target.files[0];
  if (selectedFile) {
    fileSelected = true;
  }
}


  
const fetchUserProfile = async (username) => {
    try {
        const response = await axios.get(`http://4.234.181.167:8080/users/getuserprofile/${username}`);
        const { user, userPosts, numberOfFollowers, numberOfFollowees } = response.data;
        profile = {
            userID: user.id,
            username: user.username,
            firstName: user.first_name,
            lastName: user.last_name,
            bio: user.bio || "",
            profilePic: `data:image/png;base64,${user.profile_pic}`,
            createdAt: new Date(user.created_at),
            numberOfPosts: userPosts.length,
            numberOfFollowers,
            numberOfFollowees,
        };
        posts = userPosts.map(post => post.media);
    } catch (error) {
        console.error("Error fetching user:", error);
        return null;
    }
};

const fetchFollowStatus = async (CurrentUserId, followeeId) => {
      try {
        
        const response = await fetch(`http://4.234.181.167:8080/protected/isFollowing`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${localStorage.getItem("token")}`,
          },
          body: JSON.stringify({ followerId: CurrentUserId, followeeId: followeeId}),
        });
  
        if (response.ok) {
            const data = await response.json();
            isFollowing = data.isFollowing;
        } else {
          const error = await response.json();
          console.error("Failed to fetch follow status:", error);
          return;
        }
      } catch (error) {
        console.error("Error fetching follow status:", error);
      }
    };

    const toggleFollow = async () => {
      if (isFollowing) {
        await unfollowUser(get(currentUserId), profile.userID);
      } else {
        await followUser(get(currentUserId), profile.userID);
      }
    };
  
    const followUser = async (followerId: number, followeeId: number) => {
      try {
        const payload = { followerId, followeeId };
        const response = await fetch("http://4.234.181.167:8080/protected/follow", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${localStorage.getItem("token")}`,
          },
          body: JSON.stringify(payload),
        });
  
        if (response.ok) {
            isFollowing = !isFollowing;
        } else {
            const error = await response.json();
            console.error("Failed to follow user:", error);
        }
      } 
      catch (error) {
        console.error("Error following user:", error);
      }
    };
  
    const unfollowUser = async (followerId: number, followeeId: number) => {
      try {
        const response = await fetch("http://4.234.181.167:8080/protected/unfollow", {
          method: "DELETE",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${localStorage.getItem("token")}`,
          },
          body: JSON.stringify({ followerId, followeeId }),
        });
  
        if (response.ok) {
            isFollowing = !isFollowing;
        } else {
            const error = await response.json();
            console.error("Failed to unfollow user:", error);
        }
      } 
      catch (error) {
        console.error("Error unfollowing user:", error);
      }
    };


  onMount(async () => {
    await fetchUserProfile(props.username);
    if (profile) {
    await fetchFollowStatus(get(currentUserId), profile.userID);
    console.log(get(currentUserId), profile.userID);
    console.log(isFollowing);
  }
  });


</script>

<style>
  #profileBody {
    display: flex;
    justify-content: center;
    padding: 20px;
  }

  #profileOutsideContainer {
    display: flex;
    flex-direction: column;
    max-width: 816px;
    width: 100%;
    border-radius: 8px;
    padding: 0;
  }

  #profileContainer {
    display: flex;
    width: 100%;
    gap: 40px;
    padding: 40px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    margin-bottom: 1rem;
  }

  #avatarContainer {
    flex-shrink: 0;
    display: flex;
    justify-content: center;
    align-items: center;
  }

  #profileAvatar {
    width: 150px;
    height: 150px;
    border-radius: 50%;
    object-fit: cover;
  }

  #profileContent {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  #profileTop {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 10px;
  }

  #username {
    font-size: 1.6rem;
    font-weight: bold;
    color: #333;
    margin: 0;
  }

  #actions {
    display: flex;
    gap: 10px;
  }

  #profileMiddle {
    display: flex;
    gap: 30px;
    align-items: center;
  }

  #profileMiddle h3 {
    font-size: 1.2rem;
    color: #555;
    margin: 0;
  }

  #profileMiddle strong {
    color: #333;
  }

  #profileBottom {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  #profileBottom h2 {
    font-size: 1.2rem;
    color: #333;
    margin: 0;
  }

  #profileBottom p {
    font-size: 1rem;
    color: #777;
    margin: 0;
  }

  #userPostGrid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
  }

  #postImageContainer {
    width: 264px;
    height: 264px;
    overflow: hidden;
  }

  #postImageContainer img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  #postImageContainer img:hover {
    filter: brightness(80%);
    cursor: pointer;
  }

  #postsHeader {
    text-align: center;
    font-size: 18px;
    margin-top: 1rem;
    margin-bottom: 0.5rem;
    font-weight: bold;
  }
</style>

{#if profile}
<div id="profileBody">
  <div id="profileOutsideContainer">
    <div id="profileContainer">
      <!-- Avatar -->
      <div id="avatarContainer" title="Change profile picture"> 
        <img
          src={(profile.profilePic == "data:image/png;base64,null") ? "/defaultPFP.png"  : profile.profilePic}
          alt="User Avatar"
          id="profileAvatar"
        />
        <input
          type="file"
          id="file-input"
          accept="image/jpeg,image/png,image/heic,image/heif"
          style="display: none;"
          onchange={handleFileInputChange}
      />
      </div>

      <div id="profileContent">

        <div id="profileTop">
          <h1 id="username">{profile.username}</h1>
          <div id="actions">
            {#if $currentUsername == profile.username}
              <Button class="bg-gray-300 text-gray-700 px-5 py-2 rounded-lg hover:bg-gray-400 transition-colors duration-200"
              on:click={() => setEditing(true)}>Edit Profile</Button>
              {#if $editing}
                <EditProfile bio={profile.bio} avatar={profile.profilePic} />
              {/if}
            {:else}
              <Button class="bg-blue-400 text-white px-5 py-2 rounded-lg hover:bg-blue-500 transition-colors duration-200" onclick={toggleFollow}>{isFollowing ? "Unfollow" : "Follow"}</Button>
              <Button class="bg-gray-300 text-gray-700 px-5 py-2 rounded-lg hover:bg-gray-400 transition-colors duration-200" onclick={toogleWIPmessage}>Message</Button>
            {/if}
          </div>
        </div>

        <div id="profileMiddle">
          <h3><strong>{profile.numberOfPosts}</strong> posts</h3>
          <h3><strong>{profile.numberOfFollowers}</strong> followers</h3>
          <h3><strong>{profile.numberOfFollowees}</strong> following</h3>
          {#if wipMessage}
            <span style="color: red; margin-left: 160px;">Work in progress...</span>
          {/if}
        </div>

        <div id="profileBottom">
          <h2>{profile.firstName} {profile.lastName}</h2>
          <p>{profile.bio}</p>
        </div>
      </div>
      </div>
      <h5 id="postsHeader">Posts</h5>
      {#if imageClicked}
        <span style="color: red; text-align: center; margin-bottom: 0.5rem;">Post view not implemented yet</span>
      {/if}
      <div id="userPostGrid">
        {#each posts as userPost}
          <div id="postImageContainer">
            <img onclick={toogleImageClicked} id="postImage" src={`data:image/png;base64,${userPost}`} alt="posts">
          </div>
        {/each}
      </div> 
    </div>
  </div>
 
{/if}
