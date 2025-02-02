<script lang="ts">
    import { Button } from "$lib/components/ui/button";
    import { Input } from "$lib/components/ui/input";
    import axios from "axios";
    import { navigate } from "svelte-routing";
    import { checkAuthentication } from "$lib/../store";
  
    let username = $state('');
    let password = $state('');
    const handleSubmit = async (e) => {
      e.preventDefault();
      await Login(username, password);
    };

    let loginError = $state(false);


    const Login = async (username, password) => {
    try {
      const response = await axios.post("http://4.234.181.167:8080/auth/login", {
        username,
        password,
      }, {
        headers: {
          "Content-Type": "application/json",
        },
      });

      console.log("Login successful:", response.data);
      
      localStorage.setItem('token', response.data.token)
      await checkAuthentication();
      navigate('/')
      
    } catch (error) {
      loginError = true;
      console.error("Login error:", error);
    }
  };                            
  </script>


<style>
    @import url('https://fonts.googleapis.com/css2?family=Doto:wght@100..900&family=Noto+Sans:ital,wght@0,100..900;1,100..900&family=Open+Sans:ital,wght@0,300..800;1,300..800&family=Playwrite+AU+SA:wght@100..400&family=Roboto+Mono:ital,wght@0,100..700;1,100..700&display=swap');

    #loginContainer {
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
    
    @media (min-width: 600px) {
        #loginContainer {
            width: 50%;
            max-width: 400px;
        }
    }
    
    #logo {
        font-family: "Playwrite AU SA", serif;
        font-size: 36px;
        margin-bottom: 20px;
        color: #262626;
    }
    
    #credentials {
        display: flex;
        flex-direction: column;
        margin-bottom: 20px;
    }
    
    #or {
        margin: 20px 0;
        position: relative;
    }
    
    #or::before, #or::after {
        content: "";
        position: absolute;
        top: 50%;
        width: 40%;
        height: 1px;
        background-color: #dbdbdb;
    }
    
    #or::before {
        left: 0;
    }
    
    #or::after {
        right: 0;
    }
    
    #resetPassword {
        color: #0095f6;
        cursor: pointer;
    }
    
    #signUp {
        margin-top: 20px;
        font-size: 16px;
    }
    
    #signUp a {
        color: #0095f6;
        text-decoration: none;
        font-weight: bold;
    }
    </style>
    

<div id="loginContainer">
    <div id="login">
      <div id="logo">Fakestagram</div>
      <div id="credentials">
        <form id="post-form" onsubmit={handleSubmit}>
          <Input bind:value={username} on:input={() => (loginError = false)} type="text" placeholder="username" required style="margin-bottom: 10px; border: 1px solid {loginError ? 'red' : '#ccc'};" />
          <Input bind:value={password} on:input={() => (loginError = false)} type="password" placeholder="password" required style="margin-bottom: 5px; border: 1px solid {loginError ? 'red' : '#ccc'};" />
          <Button type="submit" class="bg-blue-400 text-white font-bold py-2 px-4 w-full mt-5 rounded hover:bg-blue-500 transition duration-200">Log in</Button>
        </form>
      </div>
      {#if loginError}
        <span style="color: red; font-weight:bold;">User not found</span>
      {/if}
    </div>
    <div id="or">OR</div>
    <div id="resetPassword">Forgot password? (not implemented)</div>
    <div id="signUp">
      Don't have an account? <a href="/register">Sign up</a>
    </div>
</div>

