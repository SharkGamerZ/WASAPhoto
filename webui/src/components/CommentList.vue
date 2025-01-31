<script>
export default {
	props: {
		show: {
			type: Boolean,
			required: true
		},
		photoId: {
			type: Number,
			required: true
		},
		userId: {
			type: Number,
			required: true
		},
		floating: {
			type: Boolean,
			default: false
		},
		hideInput: {
			type: Boolean,
			default: false
		},
		comments: {
			type: Array,
			default: () => []
		},
		post: {
			type: Object,
			required: false,
			default: null
		},
		isOwner: {
			type: Boolean,
			default: false
		}
	},
	data() {
		return {
			users: {}, // Cache for user data
			newComment: '',
			loading: false,
			error: null,
			localComments: [], // Add this to store comments locally
			showMenu: false
		}
	},
	emits: ['close', 'comment-added', 'like', 'delete'],
	methods: {
		async fetchComments() {
			this.loading = true;
			try {
				const response = await this.$axios.get(`/users/${this.userId}/photos/${this.photoId}/comments`);
				this.localComments = response.data || [];
				// Sort comments with newest at the bottom
				this.localComments.sort((a, b) => new Date(a.timestamp) - new Date(b.timestamp));

				// Fetch usernames for all comments
				for (let comment of this.localComments) {
					await this.getUser(comment.user_id);
				}
			} catch (e) {
				this.error = e.toString();
			}
			this.loading = false;

			// Scroll to bottom after everything is loaded and rendered
			this.$nextTick(() => {
				this.scrollToBottom();
			});
		},

		async getUser(userID) {
			// Check if user data is already cached
			if (this.users[userID]) {
				return this.users[userID];
			}

			try {
				const response = await this.$axios.get(`/users/${userID}`);
				this.users[userID] = {
					username: response.data.username,
					propic: response.data.propic
				};
				return this.users[userID];
			} catch (e) {
				console.error('Error fetching user data:', e);
				return {
					username: 'Unknown User',
					propic: null
				};
			}
		},

		async addComment() {
			if (!this.newComment.trim()) return;

			try {
				const response = await this.$axios.post(
					`/users/${this.userId}/photos/${this.photoId}/comments`,
					JSON.stringify(this.newComment.trim())
				);
				
				const newComment = {
					id: response.data.id,
					user_id: localStorage.getItem('token'),
					content: this.newComment.trim(),
					timestamp: new Date().toISOString()
				};
				
				// Add to local comments
				this.localComments.push(newComment);
				
				// Emit the new comment to parent
				this.$emit('comment-added', newComment);
				
				this.newComment = '';
				
				// Get user data for the new comment if needed
				await this.getUser(newComment.user_id);
				
				// Scroll to bottom after adding comment
				await this.$nextTick();
				this.scrollToBottom();
			} catch (e) {
				this.error = e.toString();
			}
		},

		scrollToBottom() {
			const container = this.$refs.commentsContainer;
			if (container) {
				// Force scroll to the very bottom
				container.scrollTop = container.scrollHeight;
			}
		},

		closeModal() {
			this.$emit('close');
		},

		toggleMenu() {
			this.showMenu = !this.showMenu;
		},

		handleDelete() {
			if (confirm('Are you sure you want to delete this photo?')) {
				this.$emit('delete');
			}
			this.showMenu = false;
		},

		async deleteComment(commentId) {
			try {
				await this.$axios.delete(`/users/${this.userId}/photos/${this.photoId}/comments/${commentId}`);
				this.localComments = this.localComments.filter(comment => comment.id !== commentId);
			} catch (e) {
				console.error('Error deleting comment:', e);
			}
		},

		isCommentOwner(comment) {
			return parseInt(localStorage.getItem('token')) === comment.user_id;
		}
	},
	watch: {
		comments: {
			immediate: true,
			handler(newComments) {
				if (newComments) {
					this.localComments = [...newComments];
					this.$nextTick(() => {
						this.scrollToBottom();
					});
				}
			}
		},
		show: {
			immediate: true,
			handler(newVal) {
				if (newVal) {
					this.fetchComments();
				}
			}
		}
	},
	mounted() {
		// Initial scroll to bottom
		this.$nextTick(() => {
			this.scrollToBottom();
		});
	}
}
</script>

<template>
	<div v-if="show" :class="{ 'modal-overlay': !floating }" @click.self="closeModal">
		<div :class="['comments-wrapper', { 'modal-content': !floating }]">
			<div class="comments-header">
				<div v-if="isOwner" class="menu-container">
					<button class="menu-button" @click.stop="toggleMenu">
						<svg class="feather">
							<use href="/feather-sprite-v4.29.0.svg#more-horizontal" />
						</svg>
					</button>
					<div v-if="showMenu" class="menu-dropdown">
						<button class="delete-button" @click="handleDelete">
							<svg class="feather">
								<use href="/feather-sprite-v4.29.0.svg#trash-2" />
							</svg>
							Delete Photo
						</button>
					</div>
				</div>
			</div>

			<!-- Show header only in modal mode -->
			<div v-if="!floating" class="modal-header">
				<h2>Comments</h2>
				<button class="close-button" @click="closeModal">
					<svg class="feather">
						<use href="/feather-sprite-v4.29.0.svg#x" />
					</svg>
				</button>
			</div>

			<div ref="commentsContainer" class="comments-container">
				<div v-if="loading" class="loading-spinner">
					Loading comments...
				</div>

				<div v-else-if="error" class="error-message">
					{{ error }}
				</div>

				<div v-else class="comments-list">
					<div v-if="localComments.length === 0" class="no-comments">
						No comments yet. Be the first to comment!
					</div>

					<div v-for="comment in localComments" :key="comment.id" class="comment">
						<div class="comment-header">
							<div class="user-info">
								<img v-if="users[comment.user_id]?.propic"
									:src="'data:image/png;base64,' + users[comment.user_id].propic"
									class="comment-avatar" alt="Profile Picture" />
								<strong>@{{ users[comment.user_id]?.username || 'Loading...' }}</strong>
							</div>
							<div class="comment-actions">
								<span class="timestamp">{{ new Date(comment.timestamp).toLocaleString() }}</span>
								<div v-if="isCommentOwner(comment)" class="comment-menu">
									<button class="menu-button" @click.stop="comment.showMenu = !comment.showMenu">
										<svg class="feather">
											<use href="/feather-sprite-v4.29.0.svg#more-horizontal" />
										</svg>
									</button>
									<div v-if="comment.showMenu" class="menu-dropdown">
										<button class="delete-button" @click="deleteComment(comment.id)">
											<svg class="feather">
												<use href="/feather-sprite-v4.29.0.svg#trash-2" />
											</svg>
											Delete Comment
										</button>
									</div>
								</div>
							</div>
						</div>
						<p class="comment-text">{{ comment.content }}</p>
					</div>
				</div>
			</div>

			<!-- Combined actions and input bar -->
			<div class="bottom-bar">
				<div v-if="post" class="actions-section">
					<button class="action-button" @click="$emit('like', post)">
						<svg class="feather" :class="{ 'liked': post.liked }">
							<use href="/feather-sprite-v4.29.0.svg#heart" />
						</svg>
						<span class="likes-count">{{ post.likes || 0 }} likes</span>
					</button>
				</div>

				<div v-if="!hideInput" class="comment-input">
					<input type="text" v-model="newComment" placeholder="Write a comment..." @keyup.enter="addComment" />
					<button @click="addComment" :disabled="!newComment.trim()">
						<svg class="feather">
							<use href="/feather-sprite-v4.29.0.svg#send" />
						</svg>
					</button>
				</div>
			</div>
		</div>
	</div>
</template>

<style scoped>
.modal-overlay {
	position: fixed;
	top: 0;
	left: 0;
	right: 0;
	bottom: 0;
	background: rgba(0, 0, 0, 0.5);
	display: flex;
	justify-content: center;
	align-items: center;
	z-index: 9999;
}

.comments-wrapper {
	display: flex;
	flex-direction: column;
	background: white;
	border-radius: 8px;
	position: absolute;
	top: 0;
	bottom: 0;
	left: 0;
	right: 0;
}

.modal-content {
	width: 90%;
	max-width: 500px;
	display: flex;
	flex-direction: column;
	position: relative;
	z-index: 10000;
}

/* Different styles for floating mode */
.floating .comments-container {
	height: 100%;
	max-height: none; /* Remove max-height constraint */
	overflow: hidden; /* Remove scrollbar from container when floating */
}

/* Regular modal styles */
.modal-content .comments-container {
	max-height: 400px;
	overflow-y: auto;
}

.comments-container {
	flex: 1;
	overflow-y: auto;
	padding: 16px 16px 80px 16px;
	display: flex;
	flex-direction: column;
}

.comments-list {
	display: flex;
	flex-direction: column;
	gap: 16px;
	margin-top: auto;
}

.comment {
	padding: 8px;
	border-radius: 8px;
	background: #f8f9fa;
}

.comment-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 8px;
}

.user-info {
	display: flex;
	align-items: center;
	gap: 8px;
}

.comment-avatar {
	width: 24px;
	height: 24px;
	border-radius: 50%;
	object-fit: cover;
}

.timestamp {
	font-size: 0.8rem;
	color: #666;
}

.comment-text {
	margin: 0;
	word-break: break-word;
	color: #333;
	line-height: 1.4;
}

.bottom-bar {
	position: absolute;
	bottom: 0;
	left: 0;
	right: 0;
	display: flex;
	align-items: center;
	gap: 16px;
	padding: 12px 16px;
	background: white;
	border-top: 1px solid #dbdbdb;
	border-radius: 0 0 8px 8px;
	z-index: 2;
}

.actions-section {
	display: flex;
	align-items: center;
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

.action-button .feather {
	width: 24px;
	height: 24px;
	stroke: currentColor;
}

.action-button .feather.liked {
	fill: #ed4956;
	stroke: #ed4956;
}

.likes-count {
	font-size: 0.9rem;
	color: #262626;
}

.comment-input {
	flex: 1;
	display: flex;
	gap: 8px;
}

/* Update existing comment-input styles */
.comment-input input {
	flex: 1;
	border: 1px solid #ddd;
	border-radius: 20px;
	padding: 8px 12px;
	outline: none;
}

.comment-input button {
	background: #007bff;
	color: white;
	border: none;
	border-radius: 50%;
	width: 36px;
	height: 36px;
	display: flex;
	align-items: center;
	justify-content: center;
	cursor: pointer;
	transition: background-color 0.2s;
}

.comment-input button:disabled {
	background: #ccc;
	cursor: not-allowed;
	opacity: 0.7;
}

.comment-input button:hover:not(:disabled) {
	background: #0056b3;
}

.feather {
	width: 20px;
	height: 20px;
	stroke: currentColor;
}

.no-comments {
	text-align: center;
	color: #666;
	padding: 20px;
}

.loading-spinner {
	text-align: center;
	padding: 20px;
	color: #666;
}

.error-message {
	color: #dc3545;
	padding: 16px;
	text-align: center;
}

/* Add styles for floating mode */
.comments-wrapper {
	display: flex;
	flex-direction: column;
	height: 100%;
}

.floating .comments-container {
	padding: 12px;
	flex: 1;
	overflow-y: auto;
}

.floating .comment-input {
	border-top: 1px solid #dbdbdb;
	padding: 12px;
}

/* Add scrollbar styling */
.comments-container::-webkit-scrollbar {
	width: 8px;
}

.comments-container::-webkit-scrollbar-track {
	background: #f1f1f1;
	border-radius: 4px;
}

.comments-container::-webkit-scrollbar-thumb {
	background: #888;
	border-radius: 4px;
}

.comments-container::-webkit-scrollbar-thumb:hover {
	background: #555;
}

.comments-header {
	position: absolute;
	top: 12px;
	right: 12px;
	z-index: 3;
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
	min-width: 150px;
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

.comment-actions {
	display: flex;
	align-items: center;
	gap: 12px;
}

.comment-menu {
	position: relative;
}

.comment-menu .menu-button {
	background: none;
	border: none;
	padding: 4px;
	cursor: pointer;
	border-radius: 50%;
	display: flex;
	align-items: center;
	justify-content: center;
	opacity: 0;
	transition: opacity 0.2s, background-color 0.2s;
}

.comment:hover .menu-button {
	opacity: 1;
}

.comment-menu .menu-button:hover {
	background-color: rgba(0, 0, 0, 0.05);
}

.comment-menu .menu-button .feather {
	width: 16px;
	height: 16px;
	stroke: #666;
}

.comment-menu .menu-dropdown {
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

.comment-menu .delete-button {
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
	font-size: 0.9rem;
}

.comment-menu .delete-button:hover {
	background-color: #dc354510;
}

.comment-menu .delete-button .feather {
	width: 14px;
	height: 14px;
}
</style>
