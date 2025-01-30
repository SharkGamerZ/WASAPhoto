<script>
export default {
	data() {
		return {
			user: null,
			errormsg: null,
			loading: false,
			username: '',
			followed: false,
			followers: 0,
			followings: 0,
			photos: [],
			isFollowing: false,
			showMenu: false,
		};
	},
	computed: {
		isOwnProfile() {
			const loggedInUserId = localStorage.getItem('token');
			return this.user && loggedInUserId && parseInt(loggedInUserId) === this.user.user_id;
		}
	},
	methods: {
		async fetchUser(userId) {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get(`/users/${userId}`);
				this.user = response.data;
				this.username = response.data.username;
				this.followed = response.data.followed;
				this.followers = response.data.followers;
				this.followings = response.data.followings;
				this.posts = response.data.post_num
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
		async toggleFollow() {
			const token = localStorage.getItem('token')
			if (this.followed) {
				await this.$axios.delete(`/users/${token}/following/${this.$route.params.id}`);
				this.followers -= 1;
			} else {
				await this.$axios.put(`/users/${token}/following/${this.$route.params.id}`);
				this.followers += 1;
			}

			this.followed = !this.followed;
		},

		async fetchPhotos(userId) {
			this.$axios.get(`/users/${userId}/photos`)
				.then(response => {
					this.photos = response.data;
				})
				.catch(error => {
					this.errormsg = error.response ? error.response.data : error.toString();
				});
		},

		logout() {
			localStorage.removeItem('token');
			this.$router.push('/session');
		},

		editProfile() {
		},

		toggleMenu() {
			this.showMenu = !this.showMenu;
		},

		closeMenu(event) {
			if (!event.target.closest('.menu-container')) {
				this.showMenu = false;
			}
		}
	},
	mounted() {
		const userId = this.$route.params.id;
		this.fetchUser(userId);
		this.fetchPhotos(userId);
		document.addEventListener('click', this.closeMenu);
	},
	watch: {
		'$route.params.id': function (newId) {
			// Checks if user is null
			if (this.user === null) {
				return
			}
			if (this.$route.fullPath.startsWith('/users/') && newId !== this.user.id) {
				this.fetchUser(newId);
				this.fetchPhotos(newId);
			}
		}
	},
	unmounted() {
		document.removeEventListener('click', this.closeMenu);
	}
};
</script>

<template>
	<div class="profile-container">
		<LoadingSpinner v-if="loading" />
		<ErrorMsg v-if="errormsg" :msg="errormsg" />

		<div v-if="user" class="user-profile">
			<div class="user-header">
				<div class="profile-main">
					<div class="profile-image-container">
						<img :src="'data:image/png;base64,' + user.propic" alt="Profile Picture" class="profile-pic" />
					</div>

					<div class="profile-info">
						<div class="profile-top">
							<h1 class="username">@{{ user.username }}</h1>
							<div class="action-buttons">
								<button v-if="!isOwnProfile" @click="toggleFollow"
									:class="['follow-button', { 'following': followed }]">
									{{ followed ? 'Following' : 'Follow' }}
								</button>

								<div v-if="isOwnProfile" class="menu-container">
									<button class="menu-dots" @click.stop="toggleMenu">
										<svg class="feather">
											<use href="/feather-sprite-v4.29.0.svg#more-horizontal" />
										</svg>
									</button>
									<div v-if="showMenu" class="menu-dropdown">
										<button @click="editProfile">
											<svg class="feather">
												<use href="/feather-sprite-v4.29.0.svg#edit-2" />
											</svg>
											Edit Profile
										</button>
										<button @click="logout" class="logout-button">
											<svg class="feather">
												<use href="/feather-sprite-v4.29.0.svg#log-out" />
											</svg>
											Logout
										</button>
									</div>
								</div>
							</div>
						</div>

						<div class="stats-container">
							<div class="stat-item">
								<span class="stat-value">{{ posts }}</span>
								<span class="stat-label">Posts</span>
							</div>
							<div class="stat-item">
								<span class="stat-value">{{ followers }}</span>
								<span class="stat-label">Followers</span>
							</div>
							<div class="stat-item">
								<span class="stat-value">{{ followings }}</span>
								<span class="stat-label">Following</span>
							</div>
						</div>

						<div class="bio-container">
							<p v-if="user.bio" class="bio-text">{{ user.bio }}</p>
							<p v-else class="bio-text empty">No bio yet</p>
						</div>
					</div>
				</div>
			</div>

			<div class="photos-section">
				<h2 class="section-title">Photos</h2>
				<div class="user-photos">
					<template v-if="photos && photos.length > 0">
						<div v-for="photo in photos" :key="photo.id" class="photo-card">
							<img :src="'data:image/jpeg;base64, ' + photo.photo" alt="User Photo" />
						</div>
					</template>
					<div v-else class="no-photos">
						<svg class="feather empty-icon">
							<use href="/feather-sprite-v4.29.0.svg#image" />
						</svg>
						<p>No photos yet</p>
					</div>
				</div>
			</div>
		</div>
	</div>
</template>

<style scoped>
.profile-container {
	max-width: 935px;
	margin: 0 auto;
	padding: 30px 20px;
}

.user-profile {
	background: white;
	border-radius: 16px;
	box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
}

.user-header {
	padding: 40px;
	border-bottom: 1px solid #f0f0f0;
}

.profile-main {
	display: flex;
	gap: 80px;
}

.profile-image-container {
	flex-shrink: 0;
}

.profile-pic {
	width: 150px;
	height: 150px;
	border-radius: 50%;
	object-fit: cover;
	border: 4px solid white;
	box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.profile-info {
	flex: 1;
}

.profile-top {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 20px;
}

.username {
	font-size: 28px;
	font-weight: 600;
	margin: 0;
}

.action-buttons {
	display: flex;
	gap: 16px;
}

.follow-button {
	padding: 8px 24px;
	border-radius: 8px;
	font-weight: 600;
	transition: all 0.2s ease;
	border: 2px solid #000;
}

.follow-button:not(.following) {
	background: #000;
	color: white;
}

.follow-button.following {
	background: white;
	color: #000;
}

.stats-container {
	display: flex;
	gap: 40px;
	margin-bottom: 24px;
}

.stat-item {
	display: flex;
	flex-direction: column;
	align-items: center;
}

.stat-value {
	font-size: 20px;
	font-weight: 600;
}

.stat-label {
	font-size: 14px;
	color: #666;
	margin-top: 4px;
}

.bio-container {
	margin-top: 20px;
}

.bio-text {
	margin: 0;
	line-height: 1.5;
	color: #333;
}

.bio-text.empty {
	color: #999;
	font-style: italic;
}

.photos-section {
	padding: 40px;
}

.section-title {
	font-size: 20px;
	font-weight: 600;
	margin-bottom: 24px;
}

.user-photos {
	display: grid;
	grid-template-columns: repeat(3, 1fr);
	gap: 28px;
}

.photo-card {
	aspect-ratio: 1;
	border-radius: 12px;
	overflow: hidden;
	box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
	transition: transform 0.2s ease;
}

.photo-card:hover {
	transform: translateY(-4px);
}

.photo-card img {
	width: 100%;
	height: 100%;
	object-fit: cover;
}

.no-photos {
	grid-column: 1 / -1;
	text-align: center;
	padding: 60px 0;
	color: #666;
}

.empty-icon {
	width: 48px;
	height: 48px;
	margin-bottom: 16px;
	color: #999;
}

.menu-container {
	position: relative;
}

.menu-dots {
	background: none;
	border: none;
	padding: 8px;
	cursor: pointer;
	border-radius: 50%;
}

.menu-dots:hover {
	background: #f5f5f5;
}

.menu-dropdown {
	position: absolute;
	right: 0;
	top: 100%;
	background: white;
	border-radius: 12px;
	box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
	min-width: 180px;
	overflow: hidden;
}

.menu-dropdown button {
	display: flex;
	align-items: center;
	gap: 12px;
	width: 100%;
	padding: 12px 16px;
	border: none;
	background: none;
	cursor: pointer;
	color: #333;
}

.menu-dropdown button:hover {
	background: #f5f5f5;
}

.menu-dropdown .logout-button {
	color: #dc3545;
}

.feather {
	width: 20px;
	height: 20px;
	stroke: currentColor;
}
</style>
