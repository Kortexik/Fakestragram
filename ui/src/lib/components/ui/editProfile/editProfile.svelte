<script lang="ts">
    import { Button } from "$lib/components/ui/button";
    import Textarea from "../textarea/textarea.svelte";
    import axios from "axios";
    import { editing, setEditing } from "$lib/../store";
    import { currentUserId } from "../pages/userstore";
    import { get } from "svelte/store"
    let props = $props();

    let avatar = $state(props.avatar);
    let bio = $state(props.bio);
    let selectedFile;
    let previewAvatar = $state(avatar);

    function clickFileInput() {
        document.getElementById('avatarInput').click();
    }

    const handleAvatarChange = (e) => {
        const file = e.target.files[0];
        if (file) {
            selectedFile = file;
            previewAvatar = URL.createObjectURL(file);
        }
    };

    const Upload = async (file: File | null, bio: string) => {
        try {
            const formData = new FormData();
            if (file) {
                formData.append("avatar", file);
            }
            formData.append("bio", bio);
            formData.append("id", get(currentUserId));

            const token = localStorage.getItem("token");
            const response = await axios.put(
                "http://4.234.181.167:8080/protected/update-profile",
                formData,
                {
                    headers: {
                        "Authorization": `Bearer ${token}`,
                        "Content-Type": "multipart/form-data"
                    },
                }
            );

            if (response.status === 200) {
                setEditing(false);
                window.location.reload();
            } else {
                console.error("Failed to update profile");
            }
        } catch (error) {
            console.error("Error updating profile:", error);
        }
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        await Upload(selectedFile || null, bio);
    };


    

</script>

<style>
    #popupContainer {
        position: fixed;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        min-width: 550px;
        max-width: 90%;
        padding: 30px;
        z-index: 1000;
        background-color: #ffffff;
        border-radius: 12px;
        box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
        display: flex;
        flex-direction: column;
        align-items: center;
    }

    #backdrop {
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background: rgba(0, 0, 0, 0.5);
        z-index: 999;
    }

    form {
        width: 100%;
        display: flex;
        flex-direction: column;
        align-items: center;
    }

    #avatarSection {
        width: 100%;
        margin-bottom: 20px;
        display: flex;
        flex-direction: column; 
        align-items: center;
        text-align: left;
    }

    .Headers {
        font-size: 1.1em;
        font-weight: bold;
        width: 100%;
        text-align: left;
        margin-bottom: 8px;
    }
    #avatarContainer {
        display: flex;
        justify-content: center;
        width: 100%;
    }

    #profileAvatar {
        width: 130px;
        height: 130px;
        border-radius: 50%;
        object-fit: cover;
        cursor: pointer;
        transition: border 0.3s ease, box-shadow 0.3s ease;
    }

    #profileAvatar:hover {
        border: 3px solid #007BFF;
        box-shadow: 0 0 12px rgba(0, 123, 255, 0.4);
    }

    #bioSection {
        width: 100%;
        margin-bottom: 10px;
    }

    #header {
        margin-bottom: 20px;
        font-size: 1.3em;
    }
</style>

<div id="backdrop"
     onclick={() => setEditing(false)}
     role="button"
     tabindex="0"
     onkeydown={(e) => {
        if (e.key === "Enter" || e.key === " ") {
            setEditing(false);
        }
     }}>
</div>

<div id="popupContainer">
    <form onsubmit={handleSubmit}>
        <span id="header">Edit your profile</span>
        <div id="avatarSection">
            <span class="Headers">Profile picture</span>
            <div id="avatarContainer" title="Change profile picture" onclick={clickFileInput} role="button" tabindex="0" onkeydown={(e) => {
                if (e.key === "Enter" || e.key === " ") {
                    clickFileInput();
                }
             }}>
                <img id="profileAvatar" src={(previewAvatar == "data:image/png;base64,null" ) ? 'defaultPFP.png' : previewAvatar} alt="profilePicture" />
                <input type="file" id="avatarInput" accept="image/*" onchange={handleAvatarChange} style="display: none;" />
            </div>
        </div>
        <div id="bioSection">
            <span class="Headers">Bio</span>
            <Textarea class="w-full mb-6 p-2 mt-3" rows={5} placeholder="Bio..." bind:value={bio} />
        </div>
        <Button type="submit" class="mt-6 px-4 py-2 w-20 bg-blue-500 text-white rounded hover:bg-blue-600">Done</Button>
    </form>
</div>