<script>
import FloatingCommentList from '@/components/FloatingCommentList.vue'
import Profiles from '@/components/Profiles.vue'

export default {
	components: {
		FloatingCommentList,
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
			showComments: false,
			showLikes: false,
			likesList: [],
			currentComments: []
		}
	},
	emits: ['like'],
	methods: {
		handleLike() {
			this.$emit('like', this.post)
		},
		navigateToUser() {
			this.$router.push(`/users/${this.userId}`)
		},
		toggleComments() {
			this.showComments = !this.showComments;
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
		}
	}
}
</script>

<template>
	<div class="photo-card">
		<div class="photo-header" @click="navigateToUser">
			<img :src="userProPic ? 'data:image/png;base64,' + userProPic : '/default-avatar.png'" alt="Profile Picture"
				class="user-avatar" />
			<span class="username">@{{ username }}</span>
		</div>

		<div class="photo-container">
			<img :src="'data:image/jpeg;base64, ' + post.photo" alt="Post Photo" />
		</div>

		<div class="photo-actions">
			<div class="action-group">
				<button class="action-button" @click="handleLike">
					<svg class="feather" :class="{ 'liked': post.liked }">
						<use href="/feather-sprite-v4.29.0.svg#heart" />
					</svg>
					Like
				</button>
				<span class="likes-count" @click="fetchLikes">
					{{ post.likes || 0 }}
				</span>
			</div>
			<button class="action-button" @click="toggleComments">
				<svg class="feather">
					<use href="/feather-sprite-v4.29.0.svg#message-circle" />
				</svg>
				Comment
			</button>
		</div>

		<div class="photo-caption">
			<p><strong>@{{ username }}</strong> {{ post.caption }}</p>
			<span class="timestamp">{{ post.timestamp }}</span>
		</div>

		<FloatingCommentList 
			:show="showComments" 
			:photo-id="post.photoID" 
			:user-id="userId"
			@close="showComments = false" 
		/>

		<Profiles :users="likesList" title="Likes" :show="showLikes" @close="showLikes = false" />
	</div>
</template>

<style scoped>
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
	cursor: pointer;
	transition: background-color 0.2s;
}

.photo-header:hover {
	background-color: #f5f5f5;
}

.user-avatar {
	width: 32px;
	height: 32px;
	border-radius: 50%;
	margin-right: 12px;
	object-fit: cover;
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

.action-group {
	display: flex;
	align-items: center;
	gap: 8px;
}

.likes-count {
	font-size: 0.9rem;
	color: #666;
	cursor: pointer;
	padding: 4px 8px;
	border-radius: 12px;
	transition: background-color 0.2s;
}

.likes-count:hover {
	background-color: #f0f0f0;
}
</style>
