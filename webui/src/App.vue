<script setup>
import { RouterLink, RouterView } from 'vue-router'
import CreatePhoto from '@/components/CreatePhoto.vue'
</script>

<script>
export default {
	name: 'App',
	components: {
		CreatePhoto
	},
	data() {
		return {
			searchQuery: '',
			searchResults: [],
			showResults: false,
			showCreatePhoto: false,
			userProPic: null,
			isSearchFocused: false,
			showLogged: false
		}
	},
	computed: {
		isLoginPage() {
			return this.$route.path === '/session';
		},
		isLoggedIn() {
			return localStorage.getItem('token') !== null;
		}
	},
	watch: {
		isLoggedIn: {
			immediate: true,
			handler(newValue) {
				if (newValue) {
					this.fetchUserProfile();
				} else {
					this.userProPic = null;
				}
			}
		},
		'$route': {
			immediate: true,
			handler(to) {
				if (to.path !== '/session' && this.isLoggedIn && !this.userProPic) {
					this.fetchUserProfile();
				}
			}
		}
	},
	methods: {
		async searchUsers() {
			if (!this.searchQuery.trim()) {
				this.searchResults = [];
				this.showResults = false;
				return;
			}

			try {
				const response = await this.$axios.get('/users', {
					params: { username: this.searchQuery }
				});
				if (response.data === null) {
					this.searchResults = [];
					this.showResults = true;
					return;
				}
				this.searchResults = response.data;
				this.showResults = true;
			} catch (e) {
				console.error('Search error:', e);
			}
		},
		navigateToUser(userId) {
			this.$router.push(`/users/${userId}`);
			this.searchQuery = '';
			this.showResults = false;
		},
		handleClickOutside(event) {
			if (!event.target.closest('.search-container')) {
				this.showResults = false;
			}
		},
		async fetchUserProfile() {
			if (this.isLoggedIn) {
				try {
					const response = await this.$axios.get(`/users/${localStorage.getItem('token')}`);
					this.userProPic = response.data.propic;
				} catch (e) {
					if (e.response && (e.response.status === 404 || e.response.status === 401)) {
						localStorage.removeItem('token');
						this.$router.push('/session');
					} else {
						console.error('Error fetching user profile:', e);
					}
				}
			}
		},
		navigateToProfile() {
			this.$router.push(`/users/${localStorage.getItem('token')}`);
		},
		handleGlobalKeydown(event) {
			if (event.target.tagName === 'INPUT' || event.target.tagName === 'TEXTAREA') {
				return;
			}

			if (event.ctrlKey || event.altKey || event.metaKey) {
				return;
			}

			if (event.key.length !== 1) {
				return;
			}

			const searchInput = this.$refs.searchInput;
			if (searchInput) {
				searchInput.focus();
				this.searchQuery = event.key;
				event.preventDefault();
			}
		},
		toggleLogged(show) {
			this.showLogged = show;
		},
		async handleLogin() {
			await this.fetchUserProfile();
		},
		handleLogout() {
			this.userProPic = null;
			this.showCreatePhoto = false;
			this.$router.push('/session');
		}
	},
	mounted() {
		document.addEventListener('click', this.handleClickOutside);
		window.addEventListener('keydown', this.handleGlobalKeydown);
		if (this.isLoggedIn) {
			this.fetchUserProfile();
		}
	},
	beforeUnmount() {
		document.removeEventListener('click', this.handleClickOutside);
		window.removeEventListener('keydown', this.handleGlobalKeydown);
	}
}
</script>

<template>
	<header v-if="!isLoginPage" class="app-header">
		<div class="header-content">
			<RouterLink to="/" class="nav-link">
				<svg class="feather">
					<use href="/feather-sprite-v4.29.0.svg#home" />
				</svg>
			</RouterLink>

			<h1 class="app-title">WASAPhoto</h1>

			<div class="search-container">
				<div class="search-bar">
					<svg class="feather search-icon">
						<use href="/feather-sprite-v4.29.0.svg#search" />
					</svg>
					<input ref="searchInput" type="text" v-model="searchQuery" @input="searchUsers"
						placeholder="Search users..." class="search-input" />
				</div>

				<div v-if="showResults" class="search-results">
					<div v-if="searchResults.length > 0">
						<div v-for="user in searchResults" :key="user.user_id" class="search-result-item"
							@click="navigateToUser(user.user_id)">
							<img :src="'data:image/png;base64,' + user.propic" alt="Profile Picture"
								class="result-profile-pic" />
							<span class="result-username">@{{ user.username }}</span>
						</div>
					</div>
					<div v-else class="no-results">
						<div class="no-results-icon">
							<svg class="feather feather-x-circle">
								<use href="/feather-sprite-v4.29.0.svg#x-circle" />
							</svg>
						</div>
						<p class="no-results-text">No users found with name "{{ searchQuery }}"</p>
					</div>
				</div>
			</div>
		</div>
	</header>

	<main :class="['main-content', { 'login-page': isLoginPage }]">
		<RouterView @login="handleLogin" @logout="handleLogout" />
	</main>

	<CreatePhoto v-if="!isLoginPage" :show="showCreatePhoto" @close="showCreatePhoto = false" />

	<footer v-if="!isLoginPage" class="app-footer">
		<button class="footer-button" @click="showCreatePhoto = true">
			<svg class="feather">
				<use href="/feather-sprite-v4.29.0.svg#plus-circle" />
			</svg>
		</button>

		<button class="footer-button profile-button" @click="navigateToProfile">
			<img v-if="userProPic" :src="'data:image/png;base64,' + userProPic" alt="Profile" class="profile-pic" />
			<svg v-else class="feather">
				<use href="/feather-sprite-v4.29.0.svg#user" />
			</svg>
		</button>
	</footer>
</template>

<style>
.app-header {
	position: fixed;
	top: 0;
	left: 0;
	right: 0;
	background: white;
	box-shadow: 0 1px 4px rgba(0, 0, 0, 0.1);
	z-index: 1000;
	height: 60px;
}

.header-content {
	display: flex;
	align-items: center;
	gap: 1rem;
	height: 100%;
	max-width: 1200px;
	margin: 0 auto;
	padding: 0 1rem;
}

.app-title {
	font-size: 1rem;
	margin: 0;
	font-weight: 600;
	color: #1a1a1a;
	white-space: nowrap;
}

.nav-link {
	display: flex;
	align-items: center;
	text-decoration: none;
	color: #333;
	padding: 0.25rem;
	border-radius: 50%;
}

.nav-link:hover {
	background-color: #f5f5f5;
}

.nav-link .feather {
	width: 18px;
	height: 18px;
}

.search-container {
	position: relative;
	flex-grow: 1;
	max-width: 400px;
}

.search-bar {
	display: flex;
	align-items: center;
	background: #f5f5f5;
	border-radius: 8px;
	padding: 0.5rem 0.75rem;
	height: 36px;
	border: 1px solid #efefef;
}

.search-icon {
	width: 16px;
	height: 16px;
	color: #666;
	margin-right: 0.5rem;
}

.search-input {
	width: 100%;
	border: none;
	background: none;
	outline: none;
	font-size: 0.9rem;
	color: #333;
	padding: 0;
}

.search-results {
	position: absolute;
	top: 64%;
	left: 0;
	right: 0;
	background: white;
	border-radius: 0 0 8px 8px;
	box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
	max-height: 400px;
	overflow-y: auto;
	border: 1px solid #efefef;
	border-top: none;
	margin-top: -1px;
}

.search-result-item:first-child {
	border-top: 1px solid #efefef;
}

.search-result-item {
	display: flex;
	align-items: center;
	padding: 0.75rem;
	cursor: pointer;
	transition: background-color 0.2s;
	border-bottom: 1px solid #efefef;
}

.search-result-item:last-child {
	border-bottom: none;
}

.search-result-item:hover {
	background-color: #f8f8f8;
}

.result-profile-pic {
	width: 32px;
	height: 32px;
	border-radius: 50%;
	object-fit: cover;
	margin-right: 0.75rem;
}

.result-username {
	color: #333;
	font-size: 0.9rem;
	font-weight: 500;
}

.main-content {
	margin-top: 60px;
	margin-bottom: 60px;
	padding: 1rem;
}

.main-content.login-page {
	margin: 0;
	padding: 0;
}

.feather {
	width: 20px;
	height: 20px;
	stroke: currentColor;
}

.no-results {
	padding: 1.5rem;
	text-align: center;
}

.no-results-icon {
	margin-bottom: 0.75rem;
}

.no-results-icon .feather {
	width: 24px;
	height: 24px;
	stroke: #999;
	opacity: 0.8;
}

.no-results-text {
	color: #666;
	font-size: 0.9rem;
	margin: 0;
}

.app-footer {
	position: fixed;
	bottom: 0;
	left: 0;
	right: 0;
	background: white;
	box-shadow: 0 -2px 8px rgba(0, 0, 0, 0.1);
	padding: 0.5rem 1rem;
	display: flex;
	justify-content: space-between;
	align-items: center;
	z-index: 1000;
	height: 60px;
}

.footer-button {
	background: none;
	border: none;
	cursor: pointer;
	padding: 0;
	border-radius: 50%;
	transition: background-color 0.2s;
	display: flex;
	align-items: center;
	justify-content: center;
	width: 40px;
	height: 40px;
	aspect-ratio: 1;
}

.footer-button:hover {
	background-color: #f5f5f5;
}

.footer-button .feather {
	width: 24px;
	height: 24px;
	stroke: #333;
}

.profile-button {
	width: 40px;
	height: 40px;
	padding: 4px;
	display: flex;
	align-items: center;
	justify-content: center;
}

.profile-button .profile-pic {
	width: 32px;
	height: 32px;
	border-radius: 50%;
	object-fit: cover;
}

body.modal-open .app-footer {
	opacity: 0;
	pointer-events: none;
}

body.modal-open {
	overflow: hidden;
}

.photo-modal-overlay {
	position: fixed;
	top: 0;
	left: 0;
	right: 0;
	bottom: 0;
	background: rgba(0, 0, 0, 0.9);
	backdrop-filter: blur(10px);
	z-index: 2000;
	display: flex;
	align-items: center;
	justify-content: center;
}
</style>
