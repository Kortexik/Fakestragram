<script lang="ts">
  import { Input } from "$lib/components/ui/input";
  import { Button } from "$lib/components/ui/button";
  import axios from "axios";
  import { navigate } from "svelte-routing";
  import { usernameExists, emailExists, setUsernameExists, setEmailExists } from "$lib/../store"
  
  let username = '';
  let password = '';
  let firstName = '';
  let lastName = '';
  let email = '';

  const checkUsernameExists = async () => {
  try {
    const response = await axios.get("http://4.234.181.167:8080/auth/check-username", {
      params: { username },
    });
    setUsernameExists(response.data.exists)
  } catch (error) {
    console.error("Error checking username:", error);
  }
};

const checkEmailExists = async () => {
  try {
    const response = await axios.get("http://4.234.181.167:8080/auth/check-email", {
      params: { email },
    });
    setEmailExists(response.data.exists);
  } catch (error) {
    console.error("Error checking email:", error);
  }
};

  const handleSubmit = async (e: Event) => {
    e.preventDefault();
    Register(username, password, firstName, lastName, email);
  };

  const Register = async (username, password, firstName, lastName, email) => {
  try {
    const response = await axios.post("http://4.234.181.167:8080/auth/register", {
      username,
      password,
      firstName,
      lastName,
      email
    }, {
      headers: {
        "Content-Type": "application/json",
      },
    });

    console.log("Register successful:", response.data);
    navigate('/login')
    
  } catch (error) {
    console.error("Register error:", error);
  }
  };



  
</script>

<style>
  @import url('https://fonts.googleapis.com/css2?family=Doto:wght@100..900&family=Noto+Sans:ital,wght@0,100..900;1,100..900&family=Open+Sans:ital,wght@0,300..800;1,300..800&family=Playwrite+AU+SA:wght@100..400&family=Roboto+Mono:ital,wght@0,100..700;1,100..700&display=swap');

  #registerContainer {
    font-family: "Noto Sans", serif;
    margin: auto;
    text-align: center;
    font-family: Arial, sans-serif;
    border: 1px solid #dbdbdb;
    padding: 20px;
    border-radius: 10px;
    background-color: #fff;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    position: fixed;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    width: 90%;
    max-width: 400px;
    height: auto;
  }

  #registerPart span {
    font-family: "Playwrite AU SA", serif;
    font-size: 2.5rem;
    margin-bottom: 20px;
    display: inline-block;
  }

  #caption {
    font-size: 14px;
    color: black;
    margin-bottom: 20px;
  }

  #or {
    margin: 15px 0;
    font-size: 14px;
    position: relative;
    text-transform: uppercase;
  }

  #or:before,
  #or:after {
    content: "";
    height: 1px;
    background-color: #dbdbdb;
    position: absolute;
    top: 50%;
    width: 40%;
  }

  #or:before {
    left: 0;
  }

  #or:after {
    right: 0;
  }

  form {
    display: flex;
    flex-direction: column;
    gap: 10px;
  }

  form p {
    font-size: 12px;
    color: black;
    margin-top: 20px;
  }

  form a {
    color: #0095f6;
    text-decoration: none;
  }

  form a:hover {
    text-decoration: underline;
  }


  #loginPart {
    margin-top: 20px;
    font-size: 14px;
  }

  #loginPart a {
    color: #0095f6;
    text-decoration: none;
  }

  #loginPart a:hover {
    text-decoration: underline;
  }
</style>

<div id="registerContainer">
  <div id="registerPart">
    <span>Fakestagram</span>
    <div id="caption">Sign up to see photos and videos from your friends.</div>
    <div id="or">OR</div>
    <form id="post-form" on:submit|preventDefault={handleSubmit}>
      <Input
        on:input={checkEmailExists}
        bind:value={email}
        type="email"
        placeholder="Email"
        required
        style="border: 1px solid {$emailExists ? 'red' : '#ccc'};"
      />
      {#if $emailExists}
        <span style="font-size: 12px; font-family: 'Noto Sans', serif; color: red; text-align: center;">There is already an account with this email.</span>
      {/if}
    
      <Input bind:value={password} type="password" placeholder="Password" required />
      <Input bind:value={firstName} type="text" placeholder="First Name" required />
      <Input bind:value={lastName} type="text" placeholder="Last Name" required />
      
      <Input
        on:input={checkUsernameExists}
        bind:value={username}
        type="text"
        placeholder="Username"
        required
        style="border: 1px solid {$usernameExists ? 'red' : '#ccc'};"
      />
      {#if $usernameExists}
        <span style="font-size: 12px; font-family: 'Noto Sans', serif; color: red; text-align: center;">This username is taken.</span>
      {/if}
    
      <p>
        By signing up, you agree to our <a href="/terms">Terms</a>. Learn how we
        collect, use and share your data in our
        <a href="/policy">Privacy Policy</a> and how we use cookies and similar
        technology in our <a href="/cookies">Cookies Policy</a>.
      </p>
    
      <Button
        type="submit"
        id="submitButton"
        class="bg-blue-400 w-full text-white font-bold py-2 px-4 rounded hover:bg-blue-500 transition duration-200"
      >
        Next
      </Button>
    </form>
    
  </div>
  <div id="loginPart">
    <span>Have an account? <a href="/login">Log in</a></span>
  </div>
</div>
