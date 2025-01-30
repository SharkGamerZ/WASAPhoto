<script>
export default {
	data() {
		return {
			users: [],
			searchQuery: '',
			errormsg: null,
			loading: false,
		};
	},
	methods: {
		async fetchUsers() {
			this.loading = true;
			this.errormsg = null;
			try {
				const response = await this.$axios.get('/users', {
					params: { username: this.searchQuery }
				});
				this.users = response.data;
			} catch (e) {
				this.errormsg = e.response ? e.response.data : e.toString();
			}
			this.loading = false;
		},
	},
	mounted() {
		this.fetchUsers(); // Fetch all users initially
	},
	props: {
		initialUsers: {
			type: Array,
			default: () => []
		}
	},
	created() {
		if (this.initialUsers.length) {
			this.users = this.initialUsers;
		}
	}
};
</script>

<template>
	<div class="search-container">
		<h1>Search Users</h1>
		<LoadingSpinner v-if="loading" />
		<ErrorMsg v-if="errormsg" :msg="errormsg" />
		<div class="search-bar">
			<label for="search" class="form-label">Search by Username</label>
			<input type="text" class="form-control" id="search" v-model="searchQuery" @input="fetchUsers"
				placeholder="Enter username" />
		</div>

		<ul class="user-list">
			<li v-for="user in users" :key="user.user_id" class="user-item"
				@click="$router.push(`/users/${user.user_id}`)">
				<img :src="'data:image/png;base64,' + user.propic" alt="Profile Picture" class="profile-pic" />
				<h5 class="username">{{ user.username }}</h5>
			</li>
		</ul>
	</div>
</template>

<style>
.search-container {
	max-width: 600px;
	margin: 0 auto;
	padding: 20px;
}

.search-bar {
	margin-bottom: 20px;
}

.user-list {
	list-style: none;
	padding: 0;
}

.user-item {
	display: flex;
	align-items: center;
	padding: 10px;
	border-bottom: 1px solid #eee;
	cursor: pointer;
	transition: background-color 0.3s;
}

.user-item:hover {
	background-color: #f9f9f9;
}

.profile-pic {
	width: 50px;
	height: 50px;
	border-radius: 50%;
	object-fit: cover;
	margin-right: 15px;
}

.username {
	font-size: 1.1em;
	color: #333;
}
</style>
