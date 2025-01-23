import { createApp, reactive } from 'vue'
import App from './App.vue'
import router from './router'
import axios from './services/axios.js';
import ErrorMsg from './components/ErrorMsg.vue'
import LoadingSpinner from './components/LoadingSpinner.vue'

import './assets/dashboard.css'
import './assets/main.css'

const app = createApp(App)

// const cors = require('cors');
/// In webui/src/main.js or a similar setup file

// Set up an Axios interceptor
axios.interceptors.request.use(config => {
	const token = localStorage.getItem('token'); // Retrieve the token from storage
	if (token) {
		config.headers.Authorization = `Bearer ${token}`; // Set the Authorization header
	}
	return config;
}, error => {
	return Promise.reject(error);
});

// Store the user ID after login
function storeUserId(userId) {
	localStorage.setItem('userId', userId);
}

// Example usage after a successful login
// storeUserId(response.data.userID);/ app.use(cors());

app.config.globalProperties.$axios = axios;
app.component("ErrorMsg", ErrorMsg);
app.component("LoadingSpinner", LoadingSpinner);
app.use(router)
app.mount('#app')
