<script>
export default {
  props: {
    show: {
      type: Boolean,
      required: true
    },
    photoId: {
      type: Number,
      required: true
    },
    userId: {
      type: Number,
      required: true
    }
  },
  data() {
    return {
      comments: [],
      users: {}, // Cache for user data
      newComment: '',
      loading: false,
      error: null
    }
  },
  emits: ['close'],
  methods: {
    async fetchComments() {
      this.loading = true;
      try {
        const response = await this.$axios.get(`/users/${this.userId}/photos/${this.photoId}/comments`);
        this.comments = response.data || [];
        // Sort comments with newest at the bottom
        this.comments.sort((a, b) => new Date(a.timestamp) - new Date(b.timestamp));

        // Fetch usernames for all comments
        for (let comment of this.comments) {
          await this.getUser(comment.user_id);
        }

        // Scroll to bottom after comments are loaded
        this.$nextTick(() => {
          this.scrollToBottom();
        });
      } catch (e) {
        this.error = e.toString();
      }
      this.loading = false;
    },

    async getUser(userID) {
      if (this.users[userID]) return this.users[userID];

      try {
        const response = await this.$axios.get(`/users/${userID}`);
        this.users[userID] = {
          username: response.data.username,
          propic: response.data.propic
        };
        return this.users[userID];
      } catch (e) {
        console.error('Error fetching user data:', e);
        return { username: 'Unknown User', propic: null };
      }
    },

    async addComment() {
      if (!this.newComment.trim()) return;

      try {
        await this.$axios.post(
          `/users/${this.userId}/photos/${this.photoId}/comments`,
          JSON.stringify(this.newComment.trim())
        );
        this.newComment = '';
        await this.fetchComments();
      } catch (e) {
        this.error = e.toString();
      }
    },

    scrollToBottom() {
      const container = this.$refs.commentsContainer;
      if (container) {
        container.scrollTop = container.scrollHeight;
      }
    },

    closeModal() {
      this.$emit('close');
    }
  },
  watch: {
    show: {
      immediate: true,
      handler(newVal) {
        if (newVal) {
          this.fetchComments();
        }
      }
    }
  }
}
</script>

<template>
  <div v-if="show" class="modal-overlay" @click.self="closeModal">
    <div class="modal-content">
      <div class="modal-header">
        <h2>Comments</h2>
        <button class="close-button" @click="closeModal">
          <svg class="feather">
            <use href="/feather-sprite-v4.29.0.svg#x" />
          </svg>
        </button>
      </div>

      <div ref="commentsContainer" class="comments-container">
        <div v-if="loading" class="loading-spinner">
          Loading comments...
        </div>

        <div v-else-if="error" class="error-message">
          {{ error }}
        </div>

        <div v-else class="comments-list">
          <div v-if="comments.length === 0" class="no-comments">
            No comments yet. Be the first to comment!
          </div>

          <div v-for="comment in comments" :key="comment.id" class="comment">
            <div class="comment-header">
              <div class="user-info">
                <img v-if="users[comment.user_id]?.propic"
                  :src="'data:image/png;base64,' + users[comment.user_id].propic"
                  class="comment-avatar" alt="Profile Picture" />
                <strong>@{{ users[comment.user_id]?.username || 'Loading...' }}</strong>
              </div>
              <span class="timestamp">{{ new Date(comment.timestamp).toLocaleString() }}</span>
            </div>
            <p class="comment-text">{{ comment.content }}</p>
          </div>
        </div>
      </div>

      <div class="comment-input">
        <input type="text" v-model="newComment" placeholder="Write a comment..." @keyup.enter="addComment" />
        <button @click="addComment" :disabled="!newComment.trim()">
          <svg class="feather">
            <use href="/feather-sprite-v4.29.0.svg#send" />
          </svg>
        </button>
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
  z-index: 9999;
}

.modal-content {
  background: white;
  border-radius: 8px;
  width: 90%;
  max-width: 500px;
  display: flex;
  flex-direction: column;
  min-height: min-content;
  max-height: 600px;
}

.modal-header {
  padding: 16px;
  border-bottom: 1px solid #eee;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.comments-container {
  flex: 1;
  overflow-y: auto;
  padding: 16px 16px 0 16px;
  min-height: 100px;
  max-height: 400px;
}

.comments-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding-bottom: 16px;
}

.comment {
  padding: 8px;
  border-radius: 8px;
  background: #f8f9fa;
}

.comment-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.comment-avatar {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  object-fit: cover;
}

.timestamp {
  font-size: 0.8rem;
  color: #666;
}

.comment-text {
  margin: 0;
  word-break: break-word;
  color: #333;
  line-height: 1.4;
}

.comment-input {
  border-top: 1px solid #eee;
  padding: 16px;
  display: flex;
  gap: 8px;
  background: white;
  border-radius: 0 0 8px 8px;
}

.comment-input input {
  flex: 1;
  border: 1px solid #ddd;
  border-radius: 20px;
  padding: 8px 12px;
  outline: none;
}

.comment-input button {
  background: #007bff;
  color: white;
  border: none;
  border-radius: 50%;
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
}

.comment-input button:disabled {
  background: #ccc;
  cursor: not-allowed;
}

.feather {
  width: 20px;
  height: 20px;
  stroke: currentColor;
}

.no-comments {
  text-align: center;
  color: #666;
  padding: 20px;
}

.loading-spinner {
  text-align: center;
  padding: 20px;
  color: #666;
}

.error-message {
  color: #dc3545;
  padding: 16px;
  text-align: center;
}
</style> 