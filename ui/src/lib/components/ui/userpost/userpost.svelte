<script lang="ts">
    import { Textarea } from "$lib/components/ui/textarea";
    import axios from 'axios';
    import { followState, setFollowState } from "./followStore";
    import { get } from "svelte/store";
    import { Link } from 'svelte-routing'; 

    let props = $props();
    let isLiked = $state();
    let postLikes = $state(0);
    let isCommentClicked = $state();  
    let isCommentsVisible = $state(false);
    const followeeId = props["user_id"];
    let postId = props["post_id"];
    let postComments = $state([]);
    let commentContent = $state("");
    let loading = $state(true);
    const CurrentUserId = props["CurrentUserId"];

    
    const fetchPostData = async () => {
      try {
        const response = await fetch(`http://4.234.181.167:8080/posts/${postId}`, {
          headers: {
            Authorization: `Bearer ${localStorage.getItem("token")}`,
          },
        });
  
        if (response.ok) {
          const data = await response.json();

          postLikes = data.data.likes.length;
          postComments = data.data.comments;
          isLiked = data.data.likes.some((like) => like.user_id == CurrentUserId);

        } else {
          console.error("Failed to fetch post data");
        }
      } catch (error) {
        console.error("Error fetching post data:", error);
      } finally {
        loading = false;
      }
    };
  
    const toggleLike = async () => {
      if (isLiked) {
        await unlikePost();
      } else {
        await likePost();
      }
    };
  
    const likePost = async () => {
      try {
        const payload = { userID: CurrentUserId, postId };
        const response = await fetch("http://4.234.181.167:8080/protected/like", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${localStorage.getItem("token")}`,
          },
          body: JSON.stringify(payload),
        });
  
        if (response.ok) {
            isLiked = true;
            postLikes += 1;
        } else {
            const error = await response.json();
            console.error("Failed to like post:", error);
        }
      } catch (error) {
            console.error("Error liking post:", error);
      }
    };
  
    const unlikePost = async () => {
      try {
            const response = await fetch(`http://4.234.181.167:8080/protected/like/${postId}`, {
            method: "DELETE",
            headers: {
                Authorization: `Bearer ${localStorage.getItem("token")}`,
            },
            });
    
            if (!response.ok) {
            const error = await response.json();
            throw new Error(error.message || "Failed to unlike post");
            }
    
            isLiked = false;
            postLikes -= 1;
        }
     catch (error) {
                console.error("Error unliking post:", error.message);
        }
    };
  
    const toggleComment = () => {
      isCommentClicked = !isCommentClicked;
    };
  
    const makeCommentsVisible = () => {
      isCommentsVisible = !isCommentsVisible;
    };
  
    const postComment = async () => {
      try {
        const payload = { userID: CurrentUserId, postId, commentContent };
        const response = await fetch("http://4.234.181.167:8080/protected/comment", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${localStorage.getItem("token")}`,
          },
          body: JSON.stringify(payload),
        });
  
        if (response.ok) {
            const data = await response.json();
            postComments = [...postComments, data];
            commentContent = "";
            isCommentClicked = false;
        } else {
            const error = await response.json();
            console.error("Failed to comment", error);
        }
      }
     catch (error) {
            console.error("Error commenting the post:", error);
      }
    };
  
    const deleteComment = async (commentID: string) => {
      try {
        const response = await fetch(`http://4.234.181.167:8080/protected/comment/${commentID}`, {
          method: "DELETE",
          headers: {
            Authorization: `Bearer ${localStorage.getItem("token")}`,
          },
        });
  
        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.message || "Failed to delete comment");
        }
  
        postComments = postComments.filter((comment) => comment.id !== commentID);
      } 
      catch (error) {
        console.error("Error deleting comment:", error.message);
      }
    };
  
    const toggleFollow = async () => {
      const isFollowing = get(followState)[followeeId];
  
      if (isFollowing) {
        await unfollowUser(CurrentUserId, followeeId);
      } else {
        await followUser(CurrentUserId, followeeId);
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
            setFollowState(followeeId, true);
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
            setFollowState(followeeId, false);
        } else {
            const error = await response.json();
            console.error("Failed to unfollow user:", error);
        }
      } 
      catch (error) {
        console.error("Error unfollowing user:", error);
      }
    };
  
    const fetchFollowStatus = async () => {
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
            setFollowState(followeeId, data.isFollowing);
        } else {
          const error = await response.json();
          console.error("Failed to fetch follow status:", error);
          return;
        }
      } catch (error) {
        console.error("Error fetching follow status:", error);
      }
    };
  
    const fetchUsername = async (user_id: number): Promise<string> => {
      try {
        const response = await axios.get(`http://4.234.181.167:8080/users/username/${user_id}`);
        const username = response.data.data;
        return username;
      } catch (error) {
        console.error("Error fetching username:", error);
        return "Unknown User";
      }
    };
  
    (async () => {
      await console.log("currentUserId:", CurrentUserId);
      await fetchPostData();
      await fetchFollowStatus();
    })();
  </script>


<style>
    #Page {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 100%;
        height: 100%;
        overflow: hidden;
        
    }

    #userPostDiv {
        width: 100%;
        max-width: 468px;
        border: 1px solid #ddd;
        border-radius: 8px;
        background: #fff;
        overflow: hidden;
        box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
        padding-bottom: 16px;
    }

    #userDetails {
        padding: 10px;
        font-size: 14px;
        font-weight: bold;
        border-bottom: 1px solid #ddd;
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    #postImage img {
        width: 100%;
        height: auto;
        max-height: 585px;
        display: block;
    }

    #buttons {
        padding: 10px;
        display: flex;
        gap: 10px;
    }

    #buttons i {
        cursor: pointer;
        font-size: 20px;
        transition: color 0.3s ease, transform 0.2s ease;
    }

    .liked {
        color: red;
    }

    #buttons i:not(.liked):hover {
        color: gray;
    }

    #likeAmount {
        padding: 0 10px;
        font-size: 14px;
        font-weight: bold;
    }

    #postDescription {
        padding: 3px 10px;
        font-size: 14px;
    }

    #comment {
        padding-left: 10px;
        margin-bottom: 2px;
        font-size: 14px;
        color: gray;
        display: flex;
        justify-content: space-between;
        align-items: center;

    }

    #ShowComments {
        padding-left: 10px;
        margin-bottom: 2px;
        font-size: 13px;
        color:rgb(58, 52, 52);
        font-weight: bold;
    }
    

    #addComment {
    display: flex;
    align-items: center;
    gap: 10px;
    border-top: 1px solid #ddd;
    padding: 10px;
}

    #commentDeleteButton {
        font-weight: bold;
        padding-right:10px;
    }

    #captionUsername {
        font-weight: bold;
        font-size: 13px;
        transition-duration: 0.3s;
    }
    #captionUsername:hover {
      font-size: 14px;
    }

    #postComment {
        color: #0095f6;
        font-weight: bold;
    }
    #username {
      font-size: 14px;
      transition-duration: 0.3s;
    }
    #username:hover {
      font-size: 15px;    
    }

    #commentUsername {
      color: #212529;
      font-size: 14px;
      transition-duration: 0.3s;
    }

    #commentUsername:hover {
      font-size: 15px;
    }

    /* Spinner */
  .loading-spinner {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100px;
  }

  .spinner {
    border: 4px solid rgba(0, 0, 0, 0.1);
    border-left-color: #3498db;
    border-radius: 50%;
    width: 40px;
    height: 40px;
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    to {
      transform: rotate(360deg);
    }
  }
</style>

<link rel="stylesheet" href="https://cdn-uicons.flaticon.com/uicons-regular-rounded/css/uicons-regular-rounded.css">

<div id="Page">
    <div id="userPostDiv">
        {#if loading}
          <div class="loading-spinner">
            <div class="spinner"></div>
          </div>
        {:else}
        <div id="userDetails">
          <div id="userInfo">
            <Link to={`/${props["user_name"]}`}><span id="username">{props["user_name"]}</span></Link>
            <span>•</span>
            <span>{props["upload_time"]}</span>
          </div>
            {#if Number(CurrentUserId) !== Number(followeeId)}
              <div role="button" 
                onclick={toggleFollow} 
                tabindex="0" 
                onkeydown={(e) => {
                  if (e.key === "Enter" || e.key === " ") {
                    toggleFollow();
                  }
                }}>
                {#if $followState[followeeId]}
                  Unfollow
                {:else}
                  Follow
                {/if}
              </div>
            {/if}
        </div>
        
        <div id="postImage">
            <img
                src={props["media"]}
                alt="Loading post..."
                onerror={() => (props["media"] = "./placeholder.png")}
            />
        </div>

        <div id="buttons">
            <i
                class="fi fi-rr-heart"
                role="button"
                tabindex="0"
                aria-label={isLiked ? "Unlike post" : "Like post"}
                onclick={toggleLike}
                onkeydown={(e) => {
                    if (e.key === "Enter" || e.key === " ") {
                        toggleLike();
                    }
                }}
                class:liked={isLiked}
                title="Like"
            ></i>
            <i
                class="fi fi-rr-comment"
                title="Comment"
                onclick={toggleComment}
                onkeydown={(e) => {
                    if (e.key === "Enter" || e.key === " ") {
                        toggleComment();
                    }
                }}
                role="button"
                tabindex="0"
            ></i>
        </div>

        <div id="likeAmount">
            {postLikes} likes
        </div>

        <div id="postDescription">
            <p><Link to={`/${props["user_name"]}`}><span id="captionUsername">{props["user_name"]}</span></Link> {props["caption"]}</p>
        </div>
            <span id="ShowComments" role="button" onclick={makeCommentsVisible} tabindex="0" onkeydown={(e) => {
                if (e.key === "Enter" || e.key === " ") {
                    makeCommentsVisible();
                }
            }} >Show all comments: {postComments.length}</span>

            {#if isCommentsVisible}
                {#each postComments as comment}
                    {#await fetchUsername(comment.user_id) then username}
                        <div id="comment">
                            <div id="commentContent">
                            <Link to={`/${username}`}><span id="commentUsername">{username}</span></Link> {comment.content}
                            </div>
                            {#if comment.user_id === CurrentUserId}
                                <div role="button" tabindex="0" id="commentDeleteButton" onclick={() => deleteComment(comment.id)} onkeydown={(e) => {
                                    if (e.key === "Enter" || e.key === " ") {
                                        deleteComment(comment.id);
                                    }
                                }}>
                                    X
                                </div>
                            {/if}
                        </div>    
                    {/await}
                {/each}
            {/if}

        {#if isCommentClicked}
            <div id="addComment">
                <Textarea bind:value={commentContent} placeholder="Add a comment..." />
                {#if commentContent.trim() !== ''}
                    <span role="button" id="postComment" onclick={postComment} tabindex="0" onkeydown={(e) => {
                        if (e.key === "Enter" || e.key === " ") {
                            postComment();
                        }
                    }}>Post</span>
                {/if}
            </div>
        {/if}
      {/if}
    </div>
</div>
