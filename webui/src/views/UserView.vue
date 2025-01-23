<script>
export default {
	data() {
		return {
			user: null,
			errormsg: null,
			loading: false,
			isFollowing: false,
		};
	},
	methods: {
		async fetchUser(userId) {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get(`/users/${userId}`);
				this.user = response.data;
				this.isFollowing = response.data.isFollowing; // Assuming the API returns this info
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
		async toggleFollow() {
			if (this.isFollowing) {
				await this.$axios.delete(`/users/${this.user.id}/following/${this.user.id}`);
			} else {
				await this.$axios.put(`/users/${this.user.id}/following/${this.user.id}`);
			}
			this.isFollowing = !this.isFollowing;
		},
	},
	mounted() {
		const userId = this.$route.params.id; // Assuming the route provides a user ID
		this.fetchUser(userId);
	},
};
</script>

<template>
	<div>
		<h1>User Details</h1>
		<LoadingSpinner v-if="loading" />
		<ErrorMsg v-if="errormsg" :msg="errormsg" />
		<div v-if="user" class="user-details">
			<img :src="user.profilePicture" alt="Profile Picture" class="profile-pic" />
			<div class="user-info">
				<p><strong>{{ user.name }}</strong></p>
				<p>{{ user.email }}</p>
				<button @click="toggleFollow">
					{{ isFollowing ? 'Unfollow' : 'Follow' }}
				</button>
			</div>
		</div>
	</div>
</template>

<style>
.user-details {
	display: flex;
	align-items: center;
}

.profile-pic {
	width: 50px;
	height: 50px;
	border-radius: 50%;
	margin-right: 10px;
}

.user-info {
	display: flex;
	flex-direction: column;
}
</style>

