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
	<div>
		<h1>Search Users</h1>
		<LoadingSpinner v-if="loading" />
		<ErrorMsg v-if="errormsg" :msg="errormsg" />
		<div class="mb-3">
			<label for="search" class="form-label">Search by Username</label>
			<input type="text" class="form-control" id="search" v-model="searchQuery" @input="fetchUsers"
				placeholder="Enter username" />
		</div>

		<ul class="list-group">
			<li v-for="user in users" :key="user.userID" class="list-group-item d-flex align-items-center"
				@click="$router.push(`/users/${user.userID}`)" style="cursor: pointer;">
				<img :src="user.propic" alt="Profile Picture" class="img-thumbnail rounded-circle" width="50"
					height="50" />
				<h5 class="ms-3">{{ user.username }}</h5>
			</li>
		</ul>
	</div>
</template>

<style>
/* Add any specific styles for the search view here */
.list-group-item {
	display: flex;
	align-items: center;
}

.list-group-item img {
	margin-right: 15px;
}
</style>
