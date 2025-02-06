<script>
import Profiles from '@/components/Profiles.vue'
import FloatingPhotoCard from '@/components/FloatingPhotoCard.vue'
import FloatingBans from '@/components/FloatingBans.vue'
import DeleteConfirmationModal from '@/components/DeleteConfirmationModal.vue';

export default {
	components: {
		Profiles,
		FloatingPhotoCard,
		FloatingBans,
		DeleteConfirmationModal,
	},
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
			showFollowers: false,
			showFollowing: false,
			followersList: [],
			followingList: [],
			showPropicUpload: false,
			selectedPhotoIndex: null,
			currentComments: [],
			isEditingUsername: false,
			isEditingBio: false,
			editedUsername: '',
			editedBio: '',
			isEditMode: false,
			notification: null,
			notificationTimeout: null,
			showBans: false,
			bannedUsers: [],
			showOthersMenu: false,
			showBanConfirm: false
		};
	},
	computed: {
		isOwnProfile() {
			const loggedInUserId = localStorage.getItem('token');
			return this.user && loggedInUserId && parseInt(loggedInUserId) === this.user.user_id;
		},
		formattedBio() {
			return this.user?.bio?.split('\\n').join('\n') || 'No bio yet';
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
			this.$emit('logout');
			localStorage.removeItem('token');
			this.$router.push('/session');
		},

		toggleMenu() {
			this.showMenu = !this.showMenu;
		},

		closeMenus(event) {
			if (!event.target.closest('.menu-container')) {
				this.showMenu = false;
				this.showOthersMenu = false;
			}
		},

		async fetchFollowers() {
			try {
				const response = await this.$axios.get(`/users/${this.$route.params.id}/followers`);
				this.followersList = response.data || [];
				this.showFollowers = true;
				document.body.classList.add('modal-open');
			} catch (e) {
				this.errormsg = e.toString();
				this.followersList = [];
			}
		},

		async fetchFollowing() {
			try {
				const response = await this.$axios.get(`/users/${this.$route.params.id}/following`);
				this.followingList = response.data || [];
				this.showFollowing = true;
				document.body.classList.add('modal-open');
			} catch (e) {
				this.errormsg = e.toString();
				this.followingList = [];
			}
		},

		closeFollowers() {
			this.showFollowers = false;
			if (!this.showFollowing) {
				document.body.classList.remove('modal-open');
			}
		},

		closeFollowing() {
			this.showFollowing = false;
			if (!this.showFollowers) {
				document.body.classList.remove('modal-open');
			}
		},

		async updatePropic(event) {
			const file = event.target.files[0];
			if (!file) return;

			const formData = new FormData();
			formData.append('propic', file);

			try {
				await this.$axios.put(`/users/${this.user.user_id}/propic`, formData, {
					headers: {
						'Content-Type': 'multipart/form-data'
					}
				});
				// Refresh user data to get new propic
				await this.fetchUser(this.user.user_id);
			} catch (e) {
				this.errormsg = e.toString();
			}
		},

		triggerPropicUpload() {
			document.getElementById('propic-upload').click();
		},

		async openPhotoModal(index) {
			this.selectedPhotoIndex = index;
			document.body.classList.add('modal-open');
			await this.fetchPhotoComments(this.photos[index].photoID);
		},
		closePhotoModal() {
			this.selectedPhotoIndex = null;
			document.body.classList.remove('modal-open');
		},
		async navigatePhoto(direction) {
			const newIndex = this.selectedPhotoIndex + direction;
			if (newIndex >= 0 && newIndex < this.photos.length) {
				this.selectedPhotoIndex = newIndex;
				// Fetch comments for the new photo
				await this.fetchPhotoComments(this.photos[newIndex].photoID);
			}
		},
		async handleLike(photo) {
			const token = localStorage.getItem('token');
			try {
				if (photo.liked) {
					await this.$axios.delete(`/users/${photo.userID}/photos/${photo.photoID}/likes/${token}`);
					photo.likes--;
				} else {
					await this.$axios.put(`/users/${photo.userID}/photos/${photo.photoID}/likes/${token}`);
					photo.likes++;
				}
				photo.liked = !photo.liked;
			} catch (e) {
				console.error('Error toggling like:', e);
			}
		},
		async fetchPhotoComments(photoId) {
			try {
				const response = await this.$axios.get(`/users/${this.user.user_id}/photos/${photoId}/comments`);
				this.currentComments = response.data || [];
			} catch (e) {
				console.error('Error fetching comments:', e);
			}
		},
		async handleComment(comment) {
			// Add the comment to the current comments immediately
			this.currentComments.push({
				id: comment.id,
				user_id: comment.user_id,
				content: comment.content,
				timestamp: comment.timestamp
			});
		},
		handlePhotoDeleted(photoID) {
			// Remove the deleted photo from the photos array
			this.photos = this.photos.filter(photo => photo.photoID !== photoID);
			// Close the photo modal
			this.closePhotoModal();
		},
		handleKeydown(event) {
			if (this.selectedPhotoIndex === null) return;

			switch (event.key) {
				case 'Escape':
					this.closePhotoModal();
					break;
				case 'ArrowLeft':
					this.navigatePhoto(-1);
					break;
				case 'ArrowRight':
					this.navigatePhoto(1);
					break;
			}

			if (event.key === 'Escape') {
				if (this.showFollowers) {
					this.closeFollowers();
				}
				if (this.showFollowing) {
					this.closeFollowing();
				}
			}
		},
		cancelEdit() {
			this.isEditingUsername = false;
			this.isEditingBio = false;
			this.isEditMode = false;
			this.editedUsername = this.username;
			this.editedBio = this.user.bio || '';
		},
		startEditingUsername() {
			if (this.isOwnProfile) {
				this.editedUsername = this.username;
				this.isEditingUsername = true;
				this.$nextTick(() => {
					this.$refs.usernameInput?.focus();
				});
			}
		},
		startEditingBio() {
			if (this.isOwnProfile) {
				// Convert \n to actual newlines and limit to 4 lines
				const bioWithNewlines = (this.user.bio || '').split('\\n').join('\n');
				const bioLines = bioWithNewlines.split('\n');
				this.editedBio = bioLines.slice(0, 4).join('\n');
				this.isEditingBio = true;
				this.$nextTick(() => {
					this.$refs.bioInput?.focus();
				});
			}
		},
		showNotification(message, type = 'success') {
			// Clear any existing notification
			if (this.notificationTimeout) {
				clearTimeout(this.notificationTimeout);
			}

			this.notification = { message, type };

			// Auto-hide after 3 seconds
			this.notificationTimeout = setTimeout(() => {
				this.notification = null;
			}, 3000);
		},
		async saveUsername() {
			try {
				const response = await this.$axios.put(
					`/users/${this.user.user_id}/username`,
					JSON.stringify(this.editedUsername)
				);
				this.username = this.editedUsername;
				this.user.username = this.editedUsername;
				this.isEditingUsername = false;
				this.showNotification('Username updated successfully');
			} catch (e) {
				if (e.response?.status === 409) {
					this.showNotification('This username is already taken', 'error');
				} else if (e.response?.status === 400) {
					this.showNotification('Username is not valid', 'error');
				} else {
					this.showNotification('Failed to update username', 'error');
				}
				this.editedUsername = this.username;
			}
		},
		async saveBio() {
			if (!this.editedBio.trim()) {
				this.editedBio = '';
			}

			const bioLines = this.editedBio.split('\n');
			const limitedBio = bioLines.slice(0, 4).join('\n');

			try {
				const formattedBio = limitedBio.replace(/\r?\n/g, '\\n');
				const response = await this.$axios.put(
					`/users/${this.user.user_id}/bio`,
					JSON.stringify(formattedBio)
				);
				this.user.bio = limitedBio;
				this.isEditingBio = false;
				this.showNotification('Bio updated successfully');
			} catch (e) {
				this.showNotification('Failed to update bio', 'error');
				this.editedBio = this.user.bio || '';
			}
		},
		handleClickOutside(event) {
			if (this.isEditMode &&
				!event.target.closest('.editable') &&
				!event.target.closest('.edit-input') &&
				!event.target.closest('.menu-dropdown')) {
				// If we're editing, save before closing
				if (this.isEditingBio) {
					this.saveBio();
				} else {
					this.cancelEdit();
				}
			}
		},
		handleBioKeydown(event) {
			// If Shift+Enter, check number of newlines before allowing
			if (event.shiftKey && event.key === 'Enter') {
				const currentLines = this.editedBio.split('\n').length;
				if (currentLines >= 4) {
					event.preventDefault();
					return;
				}
				return;
			}
			// If just Enter, save the bio
			if (event.key === 'Enter') {
				event.preventDefault();
				this.saveBio();
			}
		},
		async fetchBannedUsers() {
			try {
				const response = await this.$axios.get(`/users/${this.user.user_id}/banned`);
				this.bannedUsers = response.data || [];
				this.showBans = true;
				this.showMenu = false;
				document.body.classList.add('modal-open');
			} catch (e) {
				console.error('Error fetching banned users:', e);
				this.bannedUsers = [];
			}
		},
		handleUnban(userId) {
			this.bannedUsers = this.bannedUsers.filter(user => user.user_id !== userId);
			this.showNotification('User unbanned successfully');
		},
		closeBans() {
			this.showBans = false;
			document.body.classList.remove('modal-open');
		},
		toggleOthersMenu() {
			this.showOthersMenu = !this.showOthersMenu;
		},
		confirmBanUser() {
			this.showBanConfirm = true;
		},
		cancelBanUser() {
			this.showBanConfirm = false;
		},
		async banUser() {
			try {
				const token = localStorage.getItem('token');
				await this.$axios.put(`/users/${token}/banned/${this.user.user_id}`);
				this.showOthersMenu = false;
				this.showNotification('User banned successfully');
				this.$router.push('/');
			} catch (e) {
				console.error('Error banning user:', e);
				this.showNotification('Failed to ban user', 'error');
			} finally {
				this.showBanConfirm = false;
			}
		}
	},
	mounted() {
		const userId = this.$route.params.id;
		this.fetchUser(userId);
		this.fetchPhotos(userId);
		document.addEventListener('click', this.closeMenus);
		window.addEventListener('keydown', this.handleKeydown);
		document.addEventListener('click', this.handleClickOutside);
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
		document.removeEventListener('click', this.closeMenus);
		window.removeEventListener('keydown', this.handleKeydown);
		document.removeEventListener('click', this.handleClickOutside);
		if (this.notificationTimeout) {
			clearTimeout(this.notificationTimeout);
		}
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
						<div class="propic-wrapper" :class="{ 'can-edit': isOwnProfile }">
							<img :src="'data:image/png;base64,' + user.propic" alt="Profile Picture"
								class="profile-pic" />
							<div v-if="isOwnProfile" class="edit-overlay" @click="triggerPropicUpload">
								<svg class="feather edit-icon">
									<use href="/feather-sprite-v4.29.0.svg#edit-2" />
								</svg>
							</div>
						</div>
						<input type="file" id="propic-upload" accept="image/*" @change="updatePropic"
							style="display: none" />
					</div>

					<div class="profile-info">
						<div class="profile-top">
							<div class="username" v-if="isOwnProfile">
								<div v-if="!isEditingUsername" class="editable-field" @click="startEditingUsername">
									<span>@{{ user.username }}</span>
									<svg class="feather edit-icon">
										<use href="/feather-sprite-v4.29.0.svg#edit-2" />
									</svg>
								</div>
								<input v-else type="text" v-model="editedUsername" @keyup.enter="saveUsername"
									@keyup.esc="isEditingUsername = false" @blur="isEditingUsername = false"
									ref="usernameInput" class="edit-input" />
							</div>
							<div v-else class="username">
								@{{ user.username }}
							</div>

							<div class="action-buttons">
								<button v-if="!isOwnProfile" @click="toggleFollow"
									:class="['follow-button', { 'following': followed }]">
									{{ followed ? 'Following' : 'Follow' }}
								</button>

								<!-- Menu for own profile -->
								<div v-if="isOwnProfile" class="menu-container">
									<button class="menu-button" @click.stop="toggleMenu">
										<svg class="feather">
											<use href="/feather-sprite-v4.29.0.svg#more-horizontal" />
										</svg>
									</button>
									<div v-if="showMenu" class="menu-dropdown">
										<button class="menu-item" @click="fetchBannedUsers">
											<svg class="feather">
												<use href="/feather-sprite-v4.29.0.svg#user-x" />
											</svg>
											Banned Users
										</button>
										<button class="menu-item" @click="logout">
											<svg class="feather">
												<use href="/feather-sprite-v4.29.0.svg#log-out" />
											</svg>
											Logout
										</button>
									</div>
								</div>

								<!-- Menu for other users' profiles -->
								<div v-if="!isOwnProfile" class="menu-container">
									<button class="menu-button" @click.stop="toggleOthersMenu">
										<svg class="feather">
											<use href="/feather-sprite-v4.29.0.svg#more-horizontal" />
										</svg>
									</button>
									<div v-if="showOthersMenu" class="menu-dropdown">
										<button class="menu-item warning" @click="confirmBanUser">
											<svg class="feather">
												<use href="/feather-sprite-v4.29.0.svg#user-x" />
											</svg>
											Ban User
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
							<div class="stat-item clickable" @click="fetchFollowers">
								<span class="stat-value">{{ followers }}</span>
								<span class="stat-label">Followers</span>
							</div>
							<div class="stat-item clickable" @click="fetchFollowing">
								<span class="stat-value">{{ followings }}</span>
								<span class="stat-label">Following</span>
							</div>
						</div>

						<div class="bio-container" v-if="isOwnProfile">
							<div v-if="!isEditingBio" class="bio-text editable-field" :class="{
								'empty': !user.bio
							}" @click="startEditingBio">
								<span style="white-space: pre-line">{{ formattedBio }}</span>
								<svg class="feather edit-icon">
									<use href="/feather-sprite-v4.29.0.svg#edit-2" />
								</svg>
							</div>
							<textarea v-else v-model="editedBio" class="edit-input bio-input"
								placeholder="Write your bio..." @keydown="handleBioKeydown"
								@keyup.esc="isEditingBio = false" @blur="isEditingBio = false"
								ref="bioInput"></textarea>
						</div>
						<div v-else class="bio-container">
							<div class="bio-text" :class="{ empty: !user.bio }">
								<span style="white-space: pre-line">{{ formattedBio }}</span>
							</div>
						</div>
					</div>
				</div>
			</div>

			<div class="photos-section">
				<h2 class="section-title">Photos</h2>
				<div class="user-photos">
					<template v-if="photos && photos.length > 0">
						<div v-for="(photo, index) in photos" :key="photo.id" class="photo-card"
							@click="openPhotoModal(index)">
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

		<!-- Photo Modal -->
		<div v-if="selectedPhotoIndex !== null" class="photo-modal" @click.self="closePhotoModal">
			<FloatingPhotoCard :post="photos[selectedPhotoIndex]" :username="username" :user-pro-pic="user.propic"
				:user-id="user.user_id" @like="handleLike" @navigate="navigatePhoto"
				@photo-deleted="handlePhotoDeleted" />
		</div>

		<!-- Floating Components -->
		<Profiles :users="followersList || []" title="Followers" :show="showFollowers" @close="closeFollowers" />

		<Profiles :users="followingList || []" title="Following" :show="showFollowing" @close="closeFollowing" />


		<FloatingBans :users="bannedUsers" :show="showBans" @close="closeBans" @unban="handleUnban" />

		<DeleteConfirmationModal v-if="showBanConfirm" :show="showBanConfirm" @confirm="banUser" @close="cancelBanUser"
			title="Ban User" message="Are you sure you want to ban this user? This action cannot be undone."
			acceptText="Ban" />

		<div v-if="notification" class="notification" :class="notification.type">
			{{ notification.message }}
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

.propic-wrapper {
	position: relative;
	width: 150px;
	height: 150px;
	border-radius: 50%;
}

.propic-wrapper.can-edit:hover .edit-overlay {
	opacity: 1;
}

.edit-overlay {
	position: absolute;
	top: 0;
	left: 0;
	right: 0;
	bottom: 0;
	background: rgba(0, 0, 0, 0.5);
	border-radius: 50%;
	display: flex;
	justify-content: center;
	align-items: center;
	opacity: 0;
	transition: opacity 0.2s ease;
	cursor: pointer;
}

.edit-icon {
	width: 24px;
	height: 24px;
	stroke: white;
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

.menu-button {
	background: none;
	border: none;
	padding: 8px;
	cursor: pointer;
	border-radius: 50%;
	display: flex;
	align-items: center;
	justify-content: center;
}

.menu-button:hover {
	background-color: rgba(0, 0, 0, 0.05);
}

.menu-dropdown {
	position: absolute;
	top: 100%;
	right: 0;
	background: white;
	border-radius: 8px;
	box-shadow: 0 2px 12px rgba(0, 0, 0, 0.15);
	padding: 8px;
	min-width: 180px;
	z-index: 1000;
}

.menu-item {
	display: flex;
	align-items: center;
	gap: 8px;
	width: 100%;
	padding: 8px 12px;
	border: none;
	background: none;
	color: #333;
	cursor: pointer;
	border-radius: 4px;
	font-size: 0.9rem;
	text-align: left;
}

.menu-item:hover {
	background-color: #f5f5f5;
}

.menu-item.warning {
	color: #dc3545;
}

.menu-item.warning:hover {
	background-color: #dc354510;
}

.menu-item .feather {
	width: 16px;
	height: 16px;
}

.menu-item.warning .feather {
	stroke: #dc3545;
}

.stat-item.clickable {
	cursor: pointer;
	transition: opacity 0.2s;
}

.stat-item.clickable:hover {
	opacity: 0.8;
}

.photo-modal {
	position: fixed;
	top: 0;
	left: 0;
	right: 0;
	bottom: 0;
	background: rgba(0, 0, 0, 0.8);
	backdrop-filter: blur(8px);
	display: flex;
	align-items: center;
	justify-content: center;
	z-index: 1000;
	padding: 20px;
}

.photo-modal-content {
	width: 90vw;
	max-width: 1400px;
	height: 90vh;
	position: relative;
	background: white;
	border-radius: 12px;
	overflow: hidden;
}

.nav-button {
	background: white;
	border: none;
	border-radius: 50%;
	width: 40px;
	height: 40px;
	display: flex;
	align-items: center;
	justify-content: center;
	cursor: pointer;
	position: absolute;
	top: 50%;
	transform: translateY(-50%);
	transition: all 0.2s ease;
	box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}

.nav-button:hover {
	background: #f8f9fa;
	transform: translateY(-50%) scale(1.1);
}

.nav-button.prev {
	left: 20px;
}

.nav-button.next {
	right: 20px;
}

.nav-button .feather {
	width: 24px;
	height: 24px;
	stroke: #333;
}

/* Move hover effect only to grid view photos */
.user-photos .photo-card {
	cursor: pointer;
	transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.user-photos .photo-card:hover {
	transform: translateY(-4px);
	box-shadow: 0 6px 16px rgba(0, 0, 0, 0.15);
}

/* Remove hover effects from modal */
.photo-modal-content .photo-card {
	transform: none;
	transition: none;
	box-shadow: none;
}

.photo-modal-content .photo-card:hover {
	transform: none;
	box-shadow: none;
}

/* Ensure modal scrollbar matches design */
.photo-modal-content::-webkit-scrollbar {
	width: 8px;
}

.photo-modal-content::-webkit-scrollbar-track {
	background: transparent;
}

.photo-modal-content::-webkit-scrollbar-thumb {
	background: rgba(255, 255, 255, 0.2);
	border-radius: 4px;
}

.photo-modal-content::-webkit-scrollbar-thumb:hover {
	background: rgba(255, 255, 255, 0.3);
}

.editable-field {
	display: flex;
	align-items: center;
	gap: 8px;
	cursor: pointer;
}

.editable-field .edit-icon {
	opacity: 0;
	transition: opacity 0.2s;
}

.editable-field:hover .edit-icon {
	opacity: 1;
}

.edit-input {
	padding: 8px;
	border: 2px solid #007bff;
	border-radius: 8px;
	font-size: inherit;
	font-family: inherit;
	background: white;
	outline: none;
}

.bio-input {
	width: 100%;
	min-height: 80px;
	resize: vertical;
	line-height: 1.5;
}

.empty {
	color: #666;
	font-style: italic;
}

.notification {
	position: fixed;
	top: 20px;
	right: 20px;
	padding: 12px 24px;
	border-radius: 8px;
	color: white;
	font-weight: 500;
	z-index: 1000;
	animation: slideIn 0.3s ease-out;
	box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.notification.success {
	background-color: #28a745;
}

.notification.error {
	background-color: #dc3545;
}

@keyframes slideIn {
	from {
		transform: translateX(100%);
		opacity: 0;
	}

	to {
		transform: translateX(0);
		opacity: 1;
	}
}

.ban-button {
	background: none;
	border: none;
	padding: 8px;
	cursor: pointer;
	border-radius: 50%;
	margin-right: 8px;
	display: flex;
	align-items: center;
	justify-content: center;
}

.ban-button:hover {
	background: #f5f5f5;
}

.ban-button .feather {
	width: 20px;
	height: 20px;
	stroke: #666;
}

.action-buttons {
	display: flex;
	align-items: center;
}
</style>
