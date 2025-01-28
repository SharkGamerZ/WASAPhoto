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
			return this.user && loggedInUserId && parseInt(loggedInUserId) === this.user.userID;
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
				this.posts = response.data.postNum
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

		fetchPhotos(userId) {
			this.$axios.get(`/users/${userId}/photos`)
				.then(response => {
					this.photos = response.data;
					console.log(this.photos);
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
			console.log('Edit profile clicked');
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
	<div>
		<h1>{{ username }}</h1>
		<LoadingSpinner v-if="loading" />
		<ErrorMsg v-if="errormsg" :msg="errormsg" />
		<div v-if="user" class="user-profile">
			<div class="user-header">
				<div class="user-info">
					<img :src="user.profilePicture" alt="Profile Picture" class="profile-pic" />
					<p class="username"><strong>@{{ user.username }}</strong></p>
					<button v-if="!isOwnProfile" @click="toggleFollow"
						:class="{ 'follow-button': true, 'following': followed }">
						{{ followed ? 'Unfollow' : 'Follow' }}
					</button>
				</div>

				<div v-if="isOwnProfile" class="menu-container">
					<button class="menu-dots" @click.stop="toggleMenu">
						<span></span>
						<span></span>
						<span></span>
					</button>
					<div v-if="showMenu" class="menu-dropdown">
						<button @click="editProfile">Edit Profile</button>
						<button @click="logout">Logout</button>
					</div>
				</div>

				<div class='user-bio'>
					<p v-if="user.bio">{{ user.bio }}</p>
					<p v-else>No bio yet</p>
				</div>
				<div class="user-stats">
					<span>{{ this.followers }} Followers</span>
					<span>{{ this.followings }} Following</span>
					<span>{{ this.posts }} Posts</span>
				</div>

			</div>
			<div class="user-photos">
				<template v-if="photos && photos.length > 0">
					<div v-for="photo in photos" :key="photo.id" class="photo">
						<img :src="'data:image/jpeg;base64, ' + photo.photo" alt="User Photo" />
					</div>
				</template>
				<p v-else class="no-photos">No photos yet</p>
			</div>
		</div>
	</div>
</template>

<style>
.user-profile {
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
	justify-content: space-between;
	margin-bottom: 20px;
	box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
	padding: 20px;
	background-color: #fff;
	border-radius: 10px;
	position: relative;
}

.profile-pic {
	width: 100px;
	height: 100px;
	border-radius: 50%;
	box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
}

.username {
	margin-top: 10px;
}

.user-info {
	display: flex;
	flex-direction: column;
	align-items: center;
	margin-right: 20px;
	position: relative;
}

.follow-button {
	margin-top: 10px;
	align-self: flex-start;
	margin-left: 0;
	padding: 10px 20px;
	border: none;
	border-radius: 20px;
	cursor: pointer;
	font-size: 16px;
	width: 150px;
	transition: background-color 0.3s ease-in-out, color 0.3s ease-in-out;
}

.follow-button.not-following {
	background-color: white;
	color: black;
	border: 1px solid black;
}

.follow-button.following {
	background-color: black;
	color: white;
	border: 1px solid black;
}

.user-stats {
	flex: 1;
	margin: 0;
	display: flex;
	justify-content: space-between;
	margin-top: 10px;
}

.user-stats span {
	font-size: 20px;
	display: flex;
	flex-direction: column;
	align-items: center;
}

.user-stats span::after {
	content: attr(data-label);
	font-size: 12px;
}

.user-bio {
	display: flex;
	justify-content: space-between;
	align-items: left;
	margin-bottom: 20px;
}

.user-bio p {
	text-align: left;
}

.user-photos {
	display: grid;
	grid-template-columns: repeat(4, 1fr);
	gap: 10px;
	box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
	padding: 20px;
	background-color: #fff;
	border-radius: 10px;
	margin-top: 20px;
}

.photo {
	aspect-ratio: 1;
	overflow: hidden;
}

.photo img {
	width: 100%;
	height: 100%;
	object-fit: cover;
	border-radius: 5px;
}

.no-photos {
	grid-column: 1 / -1;
	text-align: center;
	width: 100%;
}

.menu-container {
	position: relative;
	margin-left: auto;
}

.menu-dots {
	background: none;
	border: none;
	cursor: pointer;
	padding: 8px;
	display: flex;
	flex-direction: column;
	gap: 3px;
	align-items: center;
}

.menu-dots span {
	width: 4px;
	height: 4px;
	background-color: #333;
	border-radius: 50%;
	display: block;
}

.menu-dropdown {
	position: absolute;
	right: 0;
	top: 100%;
	background: white;
	border: 1px solid #ddd;
	border-radius: 4px;
	box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
	z-index: 1000;
	min-width: 150px;
}

.menu-dropdown button {
	display: block;
	width: 100%;
	padding: 8px 16px;
	text-align: left;
	background: none;
	border: none;
	cursor: pointer;
}

.menu-dropdown button:hover {
	background-color: #f5f5f5;
}

.menu-dropdown button:not(:last-child) {
	border-bottom: 1px solid #eee;
}
</style>
