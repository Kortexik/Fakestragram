<script lang="ts">
    import { Textarea } from "$lib/components/ui/textarea";

    let isLiked = $state();
    let isCommentClicked = $state();
    let props = $props();

    const toggleLike = () => {
        isLiked = !isLiked;
    };
    
    const toggleComment = () => {
        isCommentClicked = !isCommentClicked;
    };

</script>

<style>
    #Page {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 100%;
        height: 100%;
        overflow-y: hidden;
        overflow-x: hidden;
        margin-bottom: 5px;
    }       

    #userPostDiv {
        width: 100%;
        max-width: 500px;
        border: 1px solid #ddd;
        border-radius: 8px;
        background: #fff;
        overflow: hidden;
        box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    }

    #userDetails {
        padding: 10px;
        font-size: 14px;
        font-weight: bold;
        border-bottom: 1px solid #ddd;
    }

    #postImage img {
        width: 100%;
        height: auto;
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
        color:red;
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
        padding: 10px;
        font-size: 14px;
        color: gray;
    }

    #addComment {
        border-top: 1px solid #ddd;
        padding: 10px;
    }

</style>

<link rel='stylesheet' href='https://cdn-uicons.flaticon.com/uicons-regular-rounded/css/uicons-regular-rounded.css'>

<div id="Page">
    <div id="userPostDiv">
        <div id="userDetails">
            <span>{props["user_id"]} • {props["upload_time"]}</span>
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
            class:toogleComment = {isCommentClicked}
            role="button"
            tabindex="0"
            ></i>
        </div>

        <div id="likeAmount">
            {props["likes"]?.length || 0}
        </div>

        <div id="postDescription">
            {props["caption"]}
        </div>
        
        {#each props["comments"] as comment}
        <div id="comment">
            {comment.content}
        </div>
        {/each}
        
        {#if isCommentClicked}
        <div id="addComment">
                <Textarea placeholder="Add a comment..."/>
        </div>
        {/if}
    </div>
</div>
