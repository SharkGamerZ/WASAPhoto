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
		}
	},
	data() {
		return {
			users: {}, // Cache for user data
			newComment: '',
			loading: false,
			error: null,
			localComments: [] // Add this to store comments locally
		}
	},
	emits: ['close', 'comment-added'],
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
							<span class="timestamp">{{ new Date(comment.timestamp).toLocaleString() }}</span>
						</div>
						<p class="comment-text">{{ comment.content }}</p>
					</div>
				</div>
			</div>

			<!-- Only show input when not hidden -->
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
	padding: 16px 16px 0 16px;
	min-height: 100px;
	max-height: 400px;
	scroll-behavior: smooth;
	display: flex;
	flex-direction: column;
}

.comments-list {
	display: flex;
	flex-direction: column;
	gap: 16px;
	padding-bottom: 16px;
	min-height: min-content;
	margin-top: auto; /* Push content to the bottom */
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

.comment-input {
	position: sticky;
	bottom: 0;
	border-top: 1px solid #eee;
	padding: 16px;
	display: flex;
	gap: 8px;
	background: white;
	border-radius: 0 0 8px 8px;
	z-index: 1;
}

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
}

.comment-input button:disabled {
	background: #ccc;
	cursor: not-allowed;
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
</style>
