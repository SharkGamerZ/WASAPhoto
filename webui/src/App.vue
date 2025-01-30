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
			userProPic: null
		}
	},
	computed: {
		isLoggedIn() {
			return localStorage.getItem('token') !== null;
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
					console.error('Error fetching user profile:', e);
				}
			}
		},
		navigateToProfile() {
			this.$router.push(`/users/${localStorage.getItem('token')}`);
		}
	},
	mounted() {
		document.addEventListener('click', this.handleClickOutside);
		this.fetchUserProfile();
	},
	beforeUnmount() {
		document.removeEventListener('click', this.handleClickOutside);
	}
}
</script>

<template>
	<header class="app-header">
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
					<input 
						type="text" 
						v-model="searchQuery"
						@input="searchUsers"
						placeholder="Search users..."
						class="search-input"
					/>
				</div>

				<div v-if="showResults" class="search-results">
					<div v-if="searchResults.length > 0">
						<div 
							v-for="user in searchResults" 
							:key="user.user_id" 
							class="search-result-item"
							@click="navigateToUser(user.user_id)"
						>
							<img 
								:src="user.propic ? 'data:image/png;base64,' + user.propic : '/default-avatar.png'" 
								alt="Profile Picture" 
								class="result-profile-pic" 
							/>
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

	<main class="main-content">
		<RouterView />
	</main>

	<CreatePhoto 
		:show="showCreatePhoto"
		@close="showCreatePhoto = false"
	/>

	<footer v-if="isLoggedIn" class="app-footer">
		<button class="footer-button" @click="showCreatePhoto = true">
			<svg class="feather">
				<use href="/feather-sprite-v4.29.0.svg#plus-circle" />
			</svg>
		</button>

		<button class="footer-button profile-button" @click="navigateToProfile">
			<img 
				:src="userProPic ? 'data:image/png;base64,' + userProPic : '/default-avatar.png'" 
				alt="Profile" 
				class="profile-pic"
			/>
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
	height: 48px;
}

.header-content {
	display: flex;
	align-items: center;
	gap: 1rem;
	padding: 0 1rem;
	height: 100%;
	max-width: 1200px;
	margin: 0 auto;
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
	border-radius: 16px;
	padding: 0.25rem 0.75rem;
	height: 32px;
}

.search-icon {
	width: 14px;
	height: 14px;
	color: #666;
	margin-right: 0.5rem;
}

.search-input {
	width: 100%;
	border: none;
	background: none;
	outline: none;
	font-size: 0.85rem;
	color: #333;
	padding: 0;
}

.search-results {
	position: absolute;
	top: calc(100% + 4px);
	left: 0;
	right: 0;
	background: white;
	border-radius: 8px;
	box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
	max-height: 300px;
	overflow-y: auto;
}

.search-result-item {
	display: flex;
	align-items: center;
	padding: 0.5rem 0.75rem;
	cursor: pointer;
	transition: background-color 0.2s;
}

.search-result-item:hover {
	background-color: #f5f5f5;
}

.result-profile-pic {
	width: 28px;
	height: 28px;
	border-radius: 50%;
	object-fit: cover;
	margin-right: 0.75rem;
}

.result-username {
	color: #333;
	font-size: 0.85rem;
}

.main-content {
	margin-top: 48px;
	margin-bottom: 60px;
	padding: 1rem;
}

.feather {
	width: 20px;
	height: 20px;
	stroke: currentColor;
}

.no-results {
	padding: 1rem;
	text-align: center;
}

.no-results-icon {
	margin-bottom: 0.5rem;
}

.no-results-icon .feather {
	width: 24px;
	height: 24px;
	stroke: #dc3545;
	opacity: 0.8;
}

.no-results-text {
	color: #666;
	font-size: 0.85rem;
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
}

.footer-button {
	background: none;
	border: none;
	cursor: pointer;
	padding: 0.25rem;
	border-radius: 50%;
	transition: background-color 0.2s;
	display: flex;
	align-items: center;
	justify-content: center;
}

.footer-button:hover {
	background-color: #f5f5f5;
}

.footer-button .feather {
	width: 20px;
	height: 20px;
	stroke: #333;
}

.profile-button .profile-pic {
	width: 28px;
	height: 28px;
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
