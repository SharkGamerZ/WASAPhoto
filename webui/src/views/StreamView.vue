<script>

export default {
	data: function () {
		return {
			errormsg: null,
			loading: false,
			posts: [],

		}
	},
	methods: {
		async fetchPost(userID) {
			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get(`/users/${userID}/stream`);
				this.posts = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
	},
	mounted() {
		this.fetchPost(localStorage.getItem('token'))
	}
}

</script>

<template>

	<div>
		<div
			class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">Home page</h1>
			<div class="btn-toolbar mb-2 mb-md-0">
				<div class="btn-group me-2">
					<button type="button" class="btn btn-sm btn-outline-secondary" @click="refresh">
						Refresh
					</button>
					<button type="button" class="btn btn-sm btn-outline-secondary" @click="exportList">
						Export
					</button>
				</div>
				<div class="btn-group me-2">
					<button type="button" class="btn btn-sm btn-outline-primary" @click="newItem">
						New
					</button>
				</div>
			</div>
		</div>

		<div v-for="post in posts" :key="post.id">
			<h2>{{ post.title }}</h2>
			<p>{{ post.body }}</p>

			<div v-for="comment in post.comments" :key="comment.id">
				<p>{{ comment.body }}</p>

			</div>
		</div>

		<div v-if="posts === null">
			<p>No posts found</p>
		</div>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	</div>

</template>


<style></style>
