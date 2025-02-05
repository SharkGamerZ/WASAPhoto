import { createRouter, createWebHashHistory } from 'vue-router'
import UserView from '../views/UserView.vue'
import LoginView from '../views/LoginView.vue'
import StreamView from '../views/StreamView.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{
			path: '/', redirect: () => {
				const token = localStorage.getItem('token');
				if (token === 'undefined' || token === null) {
					return '/session';
				} else {
					return '/users/' + token + '/stream';
				}
			}
		},
		{ path: '/users/:id', component: UserView },
		{ path: '/users/:id/stream', component: StreamView },
		{ path: '/session', component: LoginView },
	]
})

export default router
