<script>
export default {
	data() {
		return {
			user: null,
			errormsg: null,
			loading: false,
		};
	},
	methods: {
		async fetchUser(userId) {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get(`/users/${userId}`);
				this.user = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
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
		<div v-if="user">
			<p><strong>Name:</strong> {{ user.name }}</p>
			<p><strong>Email:</strong> {{ user.email }}</p>
			<!-- Add more user details as needed -->
		</div>
	</div>
</template>

<style>
/* Add any specific styles for the user view here */
</style>
