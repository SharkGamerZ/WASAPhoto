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
};
</script>

<template>
	<div>
		<h1>Search Users</h1>
		<LoadingSpinner v-if="loading" />
		<ErrorMsg v-if="errormsg" :msg="errormsg" />
		<div class="mb-3">
			<label for="search" class="form-label">Search by Username</label>
			<input type="text" class="form-control" id="search" v-model="searchQuery" @input="fetchUsers" placeholder="Enter username" />
		</div>
		<ul class="list-group">
			<li v-for="user in users" :key="user.userID" class="list-group-item">
				<h5>{{ user.username }}</h5>
				<img :src="user.propic" alt="Profile Picture" class="img-thumbnail" width="50" height="50" />
				<p>{{ user.bio }}</p>
			</li>
		</ul>
	</div>
</template>

<style>
/* Add any specific styles for the search view here */
</style>
