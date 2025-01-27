<script>
export default {
	data() {
		return {
			user: null,
			errormsg: null,
			loading: false,
			newBio: '',
			newProfilePicture: null,
		};
	},
	methods: {
		async fetchUser() {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get('/me');
				this.user = response.data;
				this.newBio = this.user.bio;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
		async updateProfile() {
			this.loading = true;
			this.errormsg = null;
			try {
				let formData = new FormData();
				formData.append('bio', this.newBio);
				if (this.newProfilePicture) {
					formData.append('profilePicture', this.newProfilePicture);
				}
				await this.$axios.put('/me', formData);
				alert('Profile updated successfully!');
				this.fetchUser();
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
		handleFileUpload(event) {
			this.newProfilePicture = event.target.files[0];
		},
	},
	mounted() {
		this.fetchUser();
	},
};
</script>

<template>
	<div>
		<h1>My Profile</h1>
		<LoadingSpinner v-if="loading" />
		<ErrorMsg v-if="errormsg" :msg="errormsg" />
		<div v-if="user" class="user-card">
			<div class="user-header">
				<img :src="user.profilePicture" alt="Profile Picture" class="profile-pic" />
				<div class="user-stats">
					<p><strong>{{ user.name }}</strong></p>
					<p>{{ user.email }}</p>
					<div class="stats">
						<span>{{ user.followers }} Followers</span>
						<span>{{ user.following }} Following</span>
						<span>{{ user.posts }} Posts</span>
					</div>
				</div>
			</div>
			<div class="user-bio">
				<textarea v-model="newBio" placeholder="Update your bio"></textarea>
				<input type="file" @change="handleFileUpload" />
				<button @click="updateProfile">Update Profile</button>
			</div>
			<div class="user-photos">
				<div v-for="photo in user.photos" :key="photo.id" class="photo">
					<img :src="photo.url" alt="User Photo" />
				</div>
			</div>
		</div>
	</div>
</template>

<style>
.user-card {
	max-width: 600px;
	margin: 0 auto;
	padding: 20px;
	border-radius: 10px;
	box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
	background-color: #fff;
}

.user-header {
	display: flex;
	align-items: center;
	margin-bottom: 20px;
}

.profile-pic {
	width: 100px;
	height: 100px;
	border-radius: 50%;
	margin-right: 20px;
}

.user-stats {
	flex: 1;
}

.user-stats p {
	margin: 0;
}

.stats {
	display: flex;
	justify-content: space-between;
	margin-top: 10px;
}

.user-bio {
	display: flex;
	flex-direction: column;
	margin-bottom: 20px;
}

.user-bio textarea {
	width: 100%;
	height: 100px;
	margin-bottom: 10px;
}

.user-photos {
	display: grid;
	grid-template-columns: repeat(3, 1fr);
	gap: 10px;
}

.user-photos .photo img {
	width: 100%;
	border-radius: 5px;
}
</style>
