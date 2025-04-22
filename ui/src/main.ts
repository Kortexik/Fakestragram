import { mount } from 'svelte';
import './app.css';
import App from './App.svelte';

export const API_URL = 'http://localhost:8080';

const app = mount(App, {
  target: document.getElementById('app')!,
});

export default app;
