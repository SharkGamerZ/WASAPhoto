<script>
export default {
	data: function () {
		return {
			errormsg: null,
			loading: false,
			posts: [],
			usernames: {}, // Cache for usernames
		}
	},
	methods: {
		async fetchPost(userID) {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get(`/users/${userID}/stream`);
				this.posts = response.data;
				// Fetch usernames for all posts
				for (let post of this.posts) {
					await this.getUsername(post.userID);
				}
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
		async getUsername(userID) {
			// Check if username is already cached
			if (this.usernames[userID]) {
				return this.usernames[userID];
			}

			try {
				const response = await this.$axios.get(`/users/${userID}`);
				this.usernames[userID] = response.data.username;
				return response.data.username;
			} catch (e) {
				console.error('Error fetching username:', e);
				return 'Unknown User';
			}
		},
		async toggleLike(post) {
			try {
				if (post.liked) {
					await this.$axios.delete(`/users/${post.userID}/photos/${post.photoID}/likes/${localStorage.getItem('token')}`);
				} else {
					await this.$axios.put(`/users/${post.userID}/photos/${post.photoID}/likes/${localStorage.getItem('token')}`);
				}
				post.liked = !post.liked;
			} catch (e) {
				this.errormsg = e.toString();
			}
		},
	},
	mounted() {
		this.fetchPost(localStorage.getItem('token'))
	}
}
</script>

<template>
	<div class="stream-container">
		<LoadingSpinner v-if="loading" />
		<ErrorMsg v-if="errormsg" :msg="errormsg" />

		<div class="photo-feed">
			<div v-for="post in posts" :key="post.photoID" class="photo-card">
				<div class="photo-header">
					<img :src="post.propic" alt="Profile Picture" class="user-avatar" />
					<span class="username">@{{ usernames[post.userID] || 'Loading...' }}</span>
				</div>

				<div class="photo-container">
					<img :src="'data:image/jpeg;base64, ' + post.photo" alt="Post Photo" />
				</div>

				<div class="photo-actions">
					<button class="action-button" @click="toggleLike(post)">
						<svg class="feather" :class="{ 'liked': post.liked }">
							<use href="/feather-sprite-v4.29.0.svg#heart" />
						</svg>
						Like
					</button>
					<button class="action-button">
						<svg class="feather">
							<use href="/feather-sprite-v4.29.0.svg#message-circle" />
						</svg>
						Comment
					</button>
				</div>

				<div class="photo-caption">
					<p><strong>@{{ usernames[post.userID] || 'Loading...' }}</strong> {{ post.caption }}</p>
					<span class="timestamp">{{ post.timestamp }}</span>
				</div>
			</div>
		</div>

		<div v-if="posts.length === 0" class="no-posts">
			<p>No posts found</p>
		</div>
	</div>
</template>

<style scoped>
.stream-container {
	max-width: 600px;
	margin: 0 auto;
	padding: 20px;
}

.photo-feed {
	display: flex;
	flex-direction: column;
	gap: 20px;
}

.photo-card {
	background: white;
	border-radius: 8px;
	box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
	overflow: hidden;
}

.photo-header {
	display: flex;
	align-items: center;
	padding: 12px;
	border-bottom: 1px solid #eee;
}

.user-avatar {
	width: 32px;
	height: 32px;
	border-radius: 50%;
	margin-right: 12px;
}

.username {
	font-weight: 500;
	color: #333;
}

.photo-container {
	width: 100%;
	aspect-ratio: 1;
}

.photo-container img {
	width: 100%;
	height: 100%;
	object-fit: cover;
}

.photo-actions {
	display: flex;
	padding: 12px;
	gap: 16px;
	border-top: 1px solid #eee;
}

.action-button {
	display: flex;
	align-items: center;
	gap: 8px;
	background: none;
	border: none;
	color: #666;
	cursor: pointer;
	padding: 8px;
	border-radius: 4px;
}

.action-button:hover {
	background-color: #f5f5f5;
}

.action-button .feather {
	width: 20px;
	height: 20px;
	stroke: currentColor;
}

.action-button .feather.liked {
	fill: #ff4757;
	stroke: #ff4757;
}

.photo-caption {
	padding: 12px;
	color: #333;
}

.photo-caption p {
	margin: 0;
}

.timestamp {
	display: block;
	font-size: 0.8em;
	color: #666;
	margin-top: 4px;
}

.no-posts {
	text-align: center;
	padding: 40px;
	color: #666;
}
</style>
