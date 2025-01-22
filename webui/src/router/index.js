import { createRouter, createWebHashHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import UserView from '../views/UserView.vue'
import LoginView from '../views/LoginView.vue'
import SearchView from '../views/SearchView.vue' // Import the new component

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{ path: '/', component: HomeView },
		{ path: '/link1', component: HomeView },
		{ path: '/link2', component: HomeView },
		{ path: '/some/:id/link', component: HomeView },
		{ path: '/users/:id', component: UserView },
		{ path: '/users', component: SearchView },
		{ path: '/session', component: LoginView },
	]
})

export default router
