<script lang="ts">
    const handleLogout = () => {
        localStorage.removeItem("token");
        window.location.href = "/";
    };

  </script>
  
  <style>
    #logoutButton {
      display: flex;
      align-items: center;
      justify-content: center;
      padding: 10px 20px;
      background-color: whitesmoke;
      color: #333;
      border: none;
      border-radius: 5px;
      cursor: pointer;
      transition: background-color 0.3s ease, color 0.3s ease;
      box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
      position: fixed;
      bottom: 20px;
      left: 20px;
    }
  
    #logoutButton:hover {
      background-color: #e0e0e0;
      color: #000;
    }
  
    #logoutButton i {
      margin-right: 10px;
      transform: scale(1.5);
    }
  
    #logoutButton span {
      font-size: 16px;
      font-weight: bold;
    }
  </style>
  
  <link rel='stylesheet' href='https://cdn-uicons.flaticon.com/uicons-regular-rounded/css/uicons-regular-rounded.css'>
  
  <div id="logoutButton" role="button" tabindex="0" on:click={handleLogout} on:keydown={(e) => {
    if (e.key === "Enter" || e.key === " ") {
      handleLogout();
    }
  }}>
    <i class="fi fi-rr-sign-out-alt"></i>
    <span>Log out</span>
  </div>