<script>
export default {
	data() {
		return {
			username: '',
			errormsg: null,
			loading: false,
			newaccount: false,
		};
	},
	methods: {
		async login() {
			this.loading = true;
			this.errormsg = null;
			try {
				const response = await this.$axios.post('/session', { username: this.username });
				if (response.status === 201) {
					this.newaccount = true;
				}
				// Store the token and profile picture in localStorage
				const token = response.data.id;
				const profilePicture = response.data.profilePicture;
				localStorage.token = token;
				localStorage.setItem('token', token);
				localStorage.setItem('profilePicture', profilePicture);

				// Redirect to the home page
				this.$router.push('/');
			} catch (e) {
				this.errormsg = e.response ? e.response.data : e.toString();
			}
			this.loading = false;
		},
	}
};
</script>

<template>
	<div>
		<h1>Login</h1>
		<LoadingSpinner v-if="loading" />
		<ErrorMsg v-if="errormsg" :msg="errormsg" />
		<form @submit.prevent="login">
			<div class="mb-3">
				<label for="username" class="form-label">Username</label>
				<input type="text" class="form-control" id="username" v-model="username" required />
			</div>
			<button type="submit" class="btn btn-primary">Login</button>
		</form>
	</div>
</template>

<style>
/* Add any specific styles for the login view here */
</style>
