<template>
	<div class="modal-overlay" @click.self="closeModal">
		<div class="modal-content">
			<h3>{{ title }}</h3>
			<p>{{ message }}</p>
			<div class="modal-buttons">
				<button class="cancel-button" @click="closeModal">Cancel</button>
				<button class="confirm-button" @click="confirmDelete">{{ acceptText }}</button>
			</div>
		</div>
	</div>
</template>

<script>
export default {
	props: {
		show: {
			type: Boolean,
			required: true
		},
		title: {
			type: String,
			default: 'Confirm Action'
		},
		message: {
			type: String,
			default: 'Are you sure you want to proceed? This action cannot be undone.'
		},
		acceptText: {
			type: String,
			default: 'Delete'
		}
	},
	emits: ['confirm', 'close'],
	methods: {
		closeModal() {
			this.$emit('close');
		},
		confirmDelete() {
			this.$emit('confirm');
			this.closeModal();
		}
	}
}
</script>

<style scoped>
.modal-overlay {
	position: fixed;
	top: 0;
	left: 0;
	right: 0;
	bottom: 0;
	background: rgba(0, 0, 0, 0.5);
	display: flex;
	align-items: center;
	justify-content: center;
	z-index: 2000;
}

.modal-content {
	background: white;
	border-radius: 12px;
	padding: 24px;
	max-width: 400px;
	width: 90%;
	box-shadow: 0 4px 24px rgba(0, 0, 0, 0.2);
}

.modal-content h3 {
	margin: 0 0 16px 0;
	color: #333;
	font-size: 1.25rem;
}

.modal-content p {
	margin: 0 0 24px 0;
	color: #666;
	line-height: 1.5;
}

.modal-buttons {
	display: flex;
	gap: 12px;
	justify-content: flex-end;
}

.modal-buttons button {
	padding: 8px 16px;
	border-radius: 6px;
	border: none;
	font-weight: 500;
	cursor: pointer;
	transition: background-color 0.2s;
}

.cancel-button {
	background: #f0f0f0;
	color: #333;
}

.cancel-button:hover {
	background: #e0e0e0;
}

.confirm-button {
	background: #dc3545;
	color: white;
}

.confirm-button:hover {
	background: #c82333;
}
</style>
