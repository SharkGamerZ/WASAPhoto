<script>
import CommentList from '@/components/CommentList.vue'
import Profiles from '@/components/Profiles.vue'

export default {
	components: {
		CommentList,
		Profiles
	},
	props: {
		post: {
			type: Object,
			required: true
		},
		username: {
			type: String,
			required: true
		},
		userProPic: {
			type: String,
			required: false,
			default: null
		},
		userId: {
			type: Number,
			required: true
		}
	},
	data() {
		return {
			showLikes: false,
			likesList: [],
			currentComments: [],
			showMenu: false,
			showDeleteConfirm: false,
			isDeleting: false
		}
	},
	emits: ['like', 'photo-deleted'],
	computed: {
		isOwnPost() {
			return parseInt(localStorage.getItem('token')) === this.userId;
		}
	},
	watch: {
		'post.photoID': {
			immediate: true,
			handler(newPhotoId) {
				this.fetchComments();
			}
		}
	},
	methods: {
		handleLike() {
			this.$emit('like', this.post)
		},
		navigateToUser() {
			this.$router.push(`/users/${this.userId}`)
		},
		async fetchLikes() {
			try {
				const response = await this.$axios.get(`/users/${this.userId}/photos/${this.post.photoID}/likes`)
				this.likesList = response.data || []
				this.showLikes = true
			} catch (e) {
				console.error('Error fetching likes:', e)
				this.likesList = []
			}
		},
		toggleMenu() {
			this.showMenu = !this.showMenu;
		},
		async confirmDelete() {
			this.showDeleteConfirm = false;
			this.showMenu = false;
			this.isDeleting = true;

			// Wait for animation
			setTimeout(async () => {
				await this.deletePhoto();
				this.$emit('photo-deleted', this.post.photoID);
			}, 300); // Match animation duration

			this.isDeleting = false;
		},
		async deletePhoto() {
			try {
				await this.$axios.delete(`/users/${this.userId}/photos/${this.post.photoID}`);
				this.$router.push(`/users/${this.userId}`);
			} catch (e) {
				console.error('Error deleting photo:', e);
			}
		},
		async fetchComments() {
			try {
				const response = await this.$axios.get(`/users/${this.userId}/photos/${this.post.photoID}/comments`);
				this.currentComments = response.data || [];
			} catch (e) {
				console.error('Error fetching comments:', e);
			}
		}
	}
}
</script>

<template>
	<div class="floating-photo-card" :class="{ 'deleting': isDeleting }">
		<div class="floating-layout">
			<!-- Left side - Photo -->
			<div class="photo-side">
				<img :src="'data:image/jpeg;base64, ' + post.photo" alt="Post Photo" class="main-photo" />
			</div>

			<!-- Right side - Details -->
			<div class="details-side">
				<div class="caption-section">
					<div class="caption-header">
						<img :src="userProPic ? 'data:image/png;base64,' + userProPic : '/default-avatar.png'"
							class="caption-avatar" @click="navigateToUser" />
						<div class="caption-content">
							<span class="caption-username" @click="navigateToUser">@{{ username }}</span>
							<p class="caption-text">{{ post.caption }}</p>
							<span class="timestamp">{{ post.timestamp }}</span>
						</div>
						<!-- Menu button -->
						<div v-if="isOwnPost" class="menu-container">
							<button class="menu-button" @click="toggleMenu">
								<svg class="feather">
									<use href="/feather-sprite-v4.29.0.svg#more-horizontal" />
								</svg>
							</button>
							<div v-if="showMenu" class="menu-dropdown">
								<button class="delete-button" @click="showDeleteConfirm = true">
									<svg class="feather">
										<use href="/feather-sprite-v4.29.0.svg#trash-2" />
									</svg>
									Delete Photo
								</button>
							</div>
						</div>
					</div>
				</div>

				<div class="comments-section">
					<CommentList :show="true" :photo-id="post.photoID" :user-id="userId" :comments="currentComments"
						:floating="true" :post="post" :key="post.photoID" @like="handleLike"
						@comment-added="comment => currentComments.push(comment)" />
				</div>
			</div>
		</div>

		<!-- Delete confirmation modal -->
		<div v-if="showDeleteConfirm" class="delete-modal-overlay" @click.self="showDeleteConfirm = false">
			<div class="delete-modal">
				<h3>Delete Photo</h3>
				<p>Are you sure you want to delete this photo? This action cannot be undone.</p>
				<div class="delete-modal-buttons">
					<button class="cancel-button" @click="showDeleteConfirm = false">Cancel</button>
					<button class="confirm-button" @click="confirmDelete">Delete</button>
				</div>
			</div>
		</div>

		<Profiles :users="likesList" title="Likes" :show="showLikes" @close="showLikes = false" />
	</div>
</template>

<style scoped>
.floating-layout {
	display: flex;
	max-height: 90vh;
	width: 90vw;
	max-width: 1200px;
	background: white;
	border-radius: 8px;
	overflow: hidden;
}

.photo-side {
	flex: 1;
	max-width: 65%;
	background: black;
	display: flex;
	align-items: center;
	justify-content: center;
}

.main-photo {
	width: 100%;
	height: 100%;
	object-fit: contain;
}

.details-side {
	flex: 1;
	max-width: 35%;
	display: flex;
	flex-direction: column;
	border-left: 1px solid #dbdbdb;
	background: white;
	height: 90vh;
	overflow: hidden;
}

.details-header {
	padding: 14px;
	display: flex;
	align-items: center;
	gap: 12px;
	border-bottom: 1px solid #dbdbdb;
	cursor: pointer;
}

.user-avatar,
.caption-avatar {
	width: 32px;
	height: 32px;
	border-radius: 50%;
	object-fit: cover;
}

.username {
	font-weight: 600;
}

.caption-section {
	padding: 14px;
	border-bottom: 1px solid #dbdbdb;
}

.caption-header {
	display: flex;
	gap: 12px;
	align-items: flex-start;
	position: relative;
	width: 100%;
}

.caption-content {
	flex: 1;
}

.caption-username {
	font-weight: 600;
	display: block;
	margin-bottom: 4px;
}

.caption-text {
	margin: 0;
	line-height: 1.4;
}

.timestamp {
	display: block;
	font-size: 0.8em;
	color: #666;
	margin-top: 8px;
}

.comments-section {
	flex: 1;
	position: relative;
	overflow: hidden;
}

.actions-section {
	padding: 14px;
	border-top: 1px solid #dbdbdb;
	background: white;
}

.action-button {
	background: none;
	border: none;
	padding: 8px;
	cursor: pointer;
	display: flex;
	align-items: center;
	gap: 6px;
}

.feather {
	width: 24px;
	height: 24px;
	stroke: currentColor;
}

.feather.liked {
	fill: #ed4956;
	stroke: #ed4956;
}

.likes-count {
	font-size: 0.9rem;
	color: #262626;
}

.menu-container {
	position: relative;
	margin-left: auto;
	flex-shrink: 0;
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

.menu-button .feather {
	width: 20px;
	height: 20px;
	stroke: #666;
}

.menu-dropdown {
	position: absolute;
	top: 100%;
	right: 0;
	background: white;
	border-radius: 8px;
	box-shadow: 0 2px 12px rgba(0, 0, 0, 0.15);
	padding: 8px;
	min-width: 150px;
	z-index: 1000;
}

.delete-button {
	display: flex;
	align-items: center;
	gap: 8px;
	width: 100%;
	padding: 8px 12px;
	border: none;
	background: none;
	color: #dc3545;
	cursor: pointer;
	border-radius: 4px;
}

.delete-button:hover {
	background-color: #dc354510;
}

.delete-button .feather {
	width: 16px;
	height: 16px;
}

.delete-modal-overlay {
	position: fixed;
	top: 0;
	left: 0;
	right: 0;
	bottom: 0;
	background: rgba(0, 0, 0, 0.5);
	display: flex;
	align-items: center;
	justify-content: center;
	z-index: 2000;
}

.delete-modal {
	background: white;
	border-radius: 12px;
	padding: 24px;
	max-width: 400px;
	width: 90%;
	box-shadow: 0 4px 24px rgba(0, 0, 0, 0.2);
}

.delete-modal h3 {
	margin: 0 0 16px 0;
	color: #333;
	font-size: 1.25rem;
}

.delete-modal p {
	margin: 0 0 24px 0;
	color: #666;
	line-height: 1.5;
}

.delete-modal-buttons {
	display: flex;
	gap: 12px;
	justify-content: flex-end;
}

.delete-modal-buttons button {
	padding: 8px 16px;
	border-radius: 6px;
	border: none;
	font-weight: 500;
	cursor: pointer;
	transition: background-color 0.2s;
}

.cancel-button {
	background: #f0f0f0;
	color: #333;
}

.cancel-button:hover {
	background: #e0e0e0;
}

.confirm-button {
	background: #dc3545;
	color: white;
}

.confirm-button:hover {
	background: #c82333;
}

.floating-photo-card {
	transition: all 0.3s ease-out;
	transform-origin: center;
}

.floating-photo-card.deleting {
	transform: scale(0.8);
	opacity: 0;
	filter: blur(10px);
}
</style>
