<script>
export default {
  props: {
    users: {
      type: Array,
      required: true
    },
    title: {
      type: String,
      required: true
    },
    show: {
      type: Boolean,
      required: true
    }
  },
  emits: ['close'],
  methods: {
    closeModal() {
      this.$emit('close')
    },
    navigateToUser(userId) {
      this.$router.push(`/users/${userId}`)
      this.closeModal()
    }
  }
}
</script>

<template>
  <div v-if="show" class="modal-overlay" @click.self="closeModal">
    <div class="modal-content">
      <div class="modal-header">
        <h2>{{ title }}</h2>
        <button class="close-button" @click="closeModal">
          <svg class="feather">
            <use href="/feather-sprite-v4.29.0.svg#x" />
          </svg>
        </button>
      </div>
      
      <div class="users-list">
        <div v-if="users.length === 0" class="no-users">
          No users to display
        </div>
        <div v-else>
          <div v-for="user in users" 
               :key="user.user_id" 
               class="user-item"
               @click="navigateToUser(user.user_id)">
            <img 
              :src="user.propic ? 'data:image/png;base64,' + user.propic : '/default-avatar.png'" 
              alt="Profile Picture" 
              class="user-avatar"
            />
            <span class="username">@{{ user.username }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: 12px;
  width: 90%;
  max-width: 500px;
  max-height: 80vh;
  overflow-y: auto;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid #eee;
}

.modal-header h2 {
  margin: 0;
  font-size: 1.25rem;
}

.close-button {
  background: none;
  border: none;
  cursor: pointer;
  padding: 8px;
}

.users-list {
  padding: 16px;
}

.user-item {
  display: flex;
  align-items: center;
  padding: 12px;
  cursor: pointer;
  transition: background-color 0.2s;
  border-radius: 8px;
}

.user-item:hover {
  background-color: #f5f5f5;
}

.user-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
  margin-right: 12px;
}

.username {
  font-size: 0.95rem;
  color: #333;
}

.feather {
  width: 20px;
  height: 20px;
  stroke: currentColor;
}

.no-users {
  text-align: center;
  padding: 16px;
  color: #666;
}
</style> 