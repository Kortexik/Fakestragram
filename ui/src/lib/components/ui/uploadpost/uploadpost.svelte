<script lang="ts">
  import { showPopup } from "$lib/../store";
  import { Button } from "$lib/components/ui/button";
  import Textarea from "../textarea/textarea.svelte";
  import axios from "axios";

  let popup;
  let selectedFile;
  let caption = "";
  let loading = false;

  function handleBackdropClick(event: MouseEvent) {
    if (!popup.contains(event.target)) {
      showPopup.set(false);
    }
  }

  function handleKeyboardClick() {
    showPopup.set(false);
  }

  function clickFileInput() {
    document.getElementById('file-input').click();
  }

  let fileSelected = false;

  function handleFileInputChange(event) {
    selectedFile = event.target.files[0];
    if (selectedFile) {
      fileSelected = true;
      updateCheckboxDisplay();
    }
  }

  function updateCheckboxDisplay() {
    const checkboxContainer = document.getElementById("file-checkbox");
    if (fileSelected) {
      checkboxContainer.innerHTML = `
        <label>
          <input type="checkbox" id="file-inserted" checked disabled />
          File successfully inserted.
        </label>`;
    } else {
      checkboxContainer.innerHTML = "";
    }
  }

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (selectedFile && caption) {
      loading = true;
      await Upload(selectedFile, caption);
      loading = false;
      showPopup.set(false);
      window.location.reload();
    } else {
      alert("Please select a file and provide a caption.");
    }
  };

  const Upload = async (file: File, caption: string) => {
    try {
      const formData = new FormData();
      formData.append("file", file);
      formData.append("caption", caption);

      const token = localStorage.getItem("token");

      const response = await axios.post("http://4.234.181.167:8080/protected/upload", formData, {
        headers: {
          "Content-Type": "multipart/form-data",
          "Authorization": `Bearer ${token}`
        },
      });

      console.log("Post uploaded successfully:", response.data);
    } catch (error) {
      console.error("Error uploading post:", error);
    }
  };
</script>

<style>
  #PopupUpload {
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    width: 30%;
    height: 60%;
    z-index: 1000;
    background-color: #ededed;
    padding: 20px;
    border-radius: 8px;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    display: flex;
    justify-content: center;
    flex-direction: column;
    align-items: center;
  }

  #Backdrop {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: rgba(0, 0, 0, 0.5);
    z-index: 999;
  }

  #file-input {
    display: none;
  }

  #post-form {
    display: flex;
    flex-direction: column;
    align-items: center;
    width: 100%;
  }

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

  span {
    font-family: 'Trebuchet MS', 'Lucida Sans Unicode', 'Lucida Grande', 'Lucida Sans', Arial, sans-serif;
    font-size: 1.3em;
  }
</style>

<link rel='stylesheet' href='https://cdn-uicons.flaticon.com/uicons-regular-rounded/css/uicons-regular-rounded.css'>

<div id="Backdrop"
     on:click={handleBackdropClick}
     role="button"
     tabindex="0"
     on:keydown={(e) => {
        if (e.key === "Enter" || e.key === " ") {
            handleKeyboardClick();
}}}>   
</div>

<div id="PopupUpload" bind:this={popup}>
  {#if loading}
    <div class="loading-spinner">
      <div class="spinner"></div>
    </div>
  {:else}
    <i class="fi fi-rr-cloud-upload-alt" style="transform: scale(6); margin: 30px;"></i>
    <form id="post-form" on:submit={handleSubmit}>
      <span>Upload photos and videos</span>
      <Button type="button" on:click={clickFileInput} style="width: 40%; margin-top: 16px;">Select from computer</Button>
      <input
        type="file"
        id="file-input"
        accept="image/jpeg,image/png,image/heic,image/heif,video/mp4,video/quicktime"
        style="display: none;"
        on:change={handleFileInputChange}
      />
      <div id="file-checkbox" style="margin-top: 10px;"></div>
      <Textarea id="caption" style="margin: 10px;" rows={5} placeholder="Write a caption... " bind:value={caption} />
      <Button type="submit" style="width: 10%;"><i class="fi fi-rr-arrow-right"></i></Button>
    </form>
  {/if}
</div>