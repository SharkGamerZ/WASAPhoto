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

				// Emit login event before redirecting
				this.$emit('login');

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
	<div class="login-container">
		<div class="login-card">
			<div class="logo-container">
				<h1 class="logo">WASAPhoto</h1>
			</div>

			<LoadingSpinner v-if="loading" />
			<ErrorMsg v-if="errormsg" :msg="errormsg" class="error-container" />

			<form @submit.prevent="login" class="login-form">
				<div class="form-group">
					<input type="text" id="username" v-model="username" placeholder="Username" required
						class="form-input" />
				</div>
				<button type="submit" class="login-button" :disabled="loading">
					{{ loading ? 'Logging in...' : 'Login' }}
				</button>
			</form>
		</div>
	</div>
</template>

<style scoped>
.login-container {
	min-height: 100vh;
	display: flex;
	align-items: center;
	justify-content: center;
	background-color: #fafafa;
	padding: 1rem;
}

.login-card {
	background: white;
	border-radius: 12px;
	padding: 2rem;
	width: 100%;
	max-width: 400px;
	box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.logo-container {
	text-align: center;
	margin-bottom: 2rem;
}

.logo {
	font-size: 2.5rem;
	font-weight: 700;
	color: #1a1a1a;
	margin: 0;
}

.login-form {
	display: flex;
	flex-direction: column;
	gap: 1rem;
}

.form-group {
	position: relative;
}

.form-input {
	width: 100%;
	padding: 0.75rem 1rem;
	border: 1px solid #e0e0e0;
	border-radius: 8px;
	font-size: 1rem;
	transition: border-color 0.2s, box-shadow 0.2s;
	outline: none;
}

.form-input:focus {
	border-color: #007bff;
	box-shadow: 0 0 0 3px rgba(0, 123, 255, 0.1);
}

.login-button {
	background: #007bff;
	color: white;
	border: none;
	border-radius: 8px;
	padding: 0.75rem 1rem;
	font-size: 1rem;
	font-weight: 500;
	cursor: pointer;
	transition: background-color 0.2s;
}

.login-button:hover:not(:disabled) {
	background: #0056b3;
}

.login-button:disabled {
	background: #ccc;
	cursor: not-allowed;
}

.error-container {
	margin-bottom: 1rem;
}
</style>
