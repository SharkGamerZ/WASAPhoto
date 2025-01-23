import { createRouter, createWebHashHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import UserView from '../views/UserView.vue'
import LoginView from '../views/LoginView.vue'
import SearchView from '../views/SearchView.vue' // Import the new component

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{ path: '/', component: HomeView },
		{ path: '/users', component: SearchView },
		{ path: '/users/:id', component: UserView },
		{ path: '/session', component: LoginView },
	]
})

export default router
