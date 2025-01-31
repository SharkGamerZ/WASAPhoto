<script>
export default {
	props: {
		show: {
			type: Boolean,
			required: true
		},
		users: {
			type: Array,
			required: true
		},
		title: {
			type: String,
			default: 'Banned Users'
		}
	},
	emits: ['close', 'unban'],
	methods: {
		async unbanUser(userId) {
			try {
				const token = localStorage.getItem('token');
				await this.$axios.delete(`/users/${token}/banned/${userId}`);
				this.$emit('unban', userId);
			} catch (e) {
				console.error('Error unbanning user:', e);
			}
		}
	}
}
</script>

<template>
	<div v-if="show" class="profiles-modal" @click.self="$emit('close')">
		<div class="profiles-content">
			<div class="profiles-header">
				<h3>{{ title }}</h3>
				<button class="close-button" @click="$emit('close')">
					<svg class="feather">
						<use href="/feather-sprite-v4.29.0.svg#x" />
					</svg>
				</button>
			</div>
			<div class="profiles-list">
				<div v-if="users.length === 0" class="no-users">
					<svg class="feather empty-icon">
						<use href="/feather-sprite-v4.29.0.svg#users" />
					</svg>
					<p>No banned users</p>
				</div>
				<div v-else v-for="user in users" :key="user.user_id" class="profile-item">
					<div class="profile-info">
						<img :src="user.propic ? 'data:image/png;base64,' + user.propic : '/default-avatar.png'"
							alt="Profile Picture" class="profile-pic" />
						<span class="username">@{{ user.username }}</span>
					</div>
					<button class="unban-button" @click="unbanUser(user.user_id)">
						<svg class="feather">
							<use href="/feather-sprite-v4.29.0.svg#user-check" />
						</svg>
						Unban
					</button>
				</div>
			</div>
		</div>
	</div>
</template>

<style scoped>
.unban-button {
	display: flex;
	align-items: center;
	gap: 8px;
	padding: 6px 12px;
	border: none;
	background: #28a745;
	color: white;
	border-radius: 6px;
	cursor: pointer;
	font-size: 0.9rem;
	transition: background-color 0.2s;
}

.unban-button:hover {
	background: #218838;
}

/* Rest of the styles are similar to Profiles.vue */
.profiles-modal {
	position: fixed;
	top: 0;
	left: 0;
	right: 0;
	bottom: 0;
	background: rgba(0, 0, 0, 0.7);
	display: flex;
	justify-content: center;
	align-items: center;
	z-index: 1000;
	backdrop-filter: blur(5px);
}

.profiles-content {
	background: white;
	border-radius: 16px;
	width: 90%;
	max-width: 400px;
	max-height: 90vh;
	overflow-y: auto;
}

.profiles-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	padding: 16px;
	border-bottom: 1px solid #eee;
}

.profiles-header h3 {
	margin: 0;
	font-size: 1.2rem;
}

.close-button {
	background: none;
	border: none;
	padding: 8px;
	cursor: pointer;
	border-radius: 50%;
}

.close-button:hover {
	background: #f5f5f5;
}

.profiles-list {
	padding: 16px;
}

.profile-item {
	display: flex;
	align-items: center;
	justify-content: space-between;
	padding: 12px;
	border-radius: 8px;
}

.profile-item:hover {
	background: #f8f9fa;
}

.profile-info {
	display: flex;
	align-items: center;
	gap: 12px;
}

.profile-pic {
	width: 40px;
	height: 40px;
	border-radius: 50%;
	object-fit: cover;
}

.username {
	font-weight: 500;
}

.no-users {
	text-align: center;
	padding: 32px;
	color: #666;
}

.empty-icon {
	width: 48px;
	height: 48px;
	stroke: #999;
	margin-bottom: 16px;
}
</style>
