<script>
import PhotoCard from '@/components/PhotoCard.vue'
import LoadingSpinner from '@/components/LoadingSpinner.vue'
import ErrorMsg from '@/components/ErrorMsg.vue'
import FloatingCommentList from '@/components/FloatingCommentList.vue'

export default {
	components: {
		PhotoCard,
		LoadingSpinner,
		ErrorMsg,
		FloatingCommentList
	},
	data: function () {
		return {
			errormsg: null,
			loading: false,
			posts: [],
			users: {}, // Cache for user data
		}
	},
	methods: {
		async fetchPost(userID) {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get(`/users/${userID}/stream`);
				this.posts = response.data || [];

				// Fetch user data for all posts
				for (let post of this.posts) {
					await this.getUser(post.userID);
				}
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
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
		async toggleLike(post) {
			try {
				if (post.liked) {
					await this.$axios.delete(`/users/${post.userID}/photos/${post.photoID}/likes/${localStorage.getItem('token')}`);
					post.likes = (post.likes || 1) - 1; // Decrement likes count
				} else {
					await this.$axios.put(`/users/${post.userID}/photos/${post.photoID}/likes/${localStorage.getItem('token')}`);
					post.likes = (post.likes || 0) + 1; // Increment likes count
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
			<PhotoCard 
				v-for="post in posts" 
				:key="post.photoID"
				:post="post"
				:username="users[post.userID]?.username || 'Loading...'"
				:userProPic="users[post.userID]?.propic"
				:userId="post.userID"
				@like="toggleLike"
			/>
		</div>

		<div v-if="posts.length === 0 && !loading" class="no-posts">
			<svg class="feather empty-icon">
				<use href="/feather-sprite-v4.29.0.svg#camera-off" />
			</svg>
			<h2 class="empty-title">Your feed is empty!</h2>
			<p class="empty-text">
				Follow your friends to see their photos here, or invite them to join WASAPhoto!
			</p>
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

.no-posts {
	text-align: center;
	padding: 40px;
	color: #666;
	display: flex;
	flex-direction: column;
	align-items: center;
	gap: 1rem;
}

.empty-icon {
	width: 48px;
	height: 48px;
	stroke: #999;
	margin-bottom: 0.5rem;
}

.empty-title {
	font-size: 1.5rem;
	color: #333;
	margin: 0;
	font-weight: 600;
}

.empty-text {
	font-size: 1rem;
	color: #666;
	margin: 0;
	max-width: 300px;
	line-height: 1.5;
}
</style>
