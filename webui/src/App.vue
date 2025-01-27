<script setup>
import { RouterLink, RouterView } from 'vue-router'
</script>
<script>
export default {
	methods: {
		async handleButton() {
			const token = localStorage.getItem('token');
			if (token === 'undefined' || token === null) {
				this.$router.push('/session');
			} else {
				this.$router.push(`/users/${token}`);
			}
		}
	},
}
</script>

<template>

	<header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
		<a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6" href="#/">WASAPhoto</a>
		<button class="navbar-toggler position-absolute d-md-none collapsed" type="button" data-bs-toggle="collapse"
			data-bs-target="#sidebarMenu" aria-controls="sidebarMenu" aria-expanded="false"
			aria-label="Toggle navigation">
			<span class="navbar-toggler-icon"></span>
		</button>
		<div class="navbar-nav ms-auto">
			<button @click="handleButton()" class="nav-link btn btn-link">
				<img :src="loginButtonIcon" alt="User Icon" class="user-icon" />
				{{ loginButtonText }}
			</button>
		</div>
	</header>

	<div class="container-fluid">
		<div class="row">
			<nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
				<div class="position-sticky pt-3 sidebar-sticky">
					<h6
						class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
						<span>Menu</span>
					</h6>
					<ul class="nav flex-column">
						<li class="nav-item">
							<RouterLink to="/" class="nav-link">
								<svg class="feather">
									<use href="/feather-sprite-v4.29.0.svg#home" />
								</svg>
								Home
							</RouterLink>
						</li>
						<li class="nav-item">
							<RouterLink to="/users/" class="nav-link">
								<svg class="feather">
									<use href="/feather-sprite-v4.29.0.svg#search" />
								</svg>
								Ricerca
							</RouterLink>
						</li>
					</ul>

				</div>
			</nav>

			<main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
				<RouterView />
			</main>
		</div>
	</div>
</template>

<style>
.user-icon {
	width: 30px;
	height: 30px;
	border-radius: 50%;
	margin-right: 10px;
}
</style>
