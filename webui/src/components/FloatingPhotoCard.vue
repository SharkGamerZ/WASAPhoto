<script>
import CommentList from '@/components/CommentList.vue'
import Profiles from '@/components/Profiles.vue'

export default {
  components: {
    CommentList,
    Profiles
  },
  props: {
    post: {
      type: Object,
      required: true
    },
    username: {
      type: String,
      required: true
    },
    userProPic: {
      type: String,
      required: false,
      default: null
    },
    userId: {
      type: Number,
      required: true
    }
  },
  data() {
    return {
      showLikes: false,
      likesList: [],
      currentComments: []
    }
  },
  emits: ['like'],
  methods: {
    handleLike() {
      this.$emit('like', this.post)
    },
    navigateToUser() {
      this.$router.push(`/users/${this.userId}`)
    },
    async fetchLikes() {
      try {
        const response = await this.$axios.get(`/users/${this.userId}/photos/${this.post.photoID}/likes`)
        this.likesList = response.data || []
        this.showLikes = true
      } catch (e) {
        console.error('Error fetching likes:', e)
        this.likesList = []
      }
    }
  }
}
</script>

<template>
  <div class="floating-photo-card">
    <div class="floating-layout">
      <!-- Left side - Photo -->
      <div class="photo-side">
        <img :src="'data:image/jpeg;base64, ' + post.photo" alt="Post Photo" class="main-photo" />
      </div>

      <!-- Right side - Details -->
      <div class="details-side">
        <div class="details-header" @click="navigateToUser">
          <img :src="userProPic ? 'data:image/png;base64,' + userProPic : '/default-avatar.png'" alt="Profile Picture"
            class="user-avatar" />
          <span class="username">@{{ username }}</span>
        </div>

        <div class="caption-section">
          <div class="caption-header">
            <img :src="userProPic ? 'data:image/png;base64,' + userProPic : '/default-avatar.png'" alt="Profile Picture"
              class="caption-avatar" />
            <div class="caption-content">
              <span class="caption-username">@{{ username }}</span>
              <p class="caption-text">{{ post.caption }}</p>
              <span class="timestamp">{{ post.timestamp }}</span>
            </div>
          </div>
        </div>

        <div class="comments-section">
          <CommentList 
            :show="true" 
            :photo-id="post.photoID" 
            :user-id="userId"
            :comments="currentComments"
            :floating="true"
            @comment-added="comment => currentComments.push(comment)"
          />
        </div>

        <div class="actions-section">
          <div class="actions-wrapper">
            <button class="action-button" @click="handleLike">
              <svg class="feather" :class="{ 'liked': post.liked }">
                <use href="/feather-sprite-v4.29.0.svg#heart" />
              </svg>
              <span class="likes-count">{{ post.likes || 0 }} likes</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <Profiles :users="likesList" title="Likes" :show="showLikes" @close="showLikes = false" />
  </div>
</template>

<style scoped>
.floating-layout {
  display: flex;
  max-height: 90vh;
  width: 90vw;
  max-width: 1200px;
  background: white;
  border-radius: 8px;
  overflow: hidden;
}

.photo-side {
  flex: 1;
  max-width: 65%;
  background: black;
  display: flex;
  align-items: center;
  justify-content: center;
}

.main-photo {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.details-side {
  flex: 1;
  max-width: 35%;
  display: flex;
  flex-direction: column;
  border-left: 1px solid #dbdbdb;
  background: white;
}

.details-header {
  padding: 14px;
  display: flex;
  align-items: center;
  gap: 12px;
  border-bottom: 1px solid #dbdbdb;
  cursor: pointer;
}

.user-avatar, .caption-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  object-fit: cover;
}

.username {
  font-weight: 600;
}

.caption-section {
  padding: 14px;
  border-bottom: 1px solid #dbdbdb;
}

.caption-header {
  display: flex;
  gap: 12px;
}

.caption-content {
  flex: 1;
}

.caption-username {
  font-weight: 600;
  display: block;
  margin-bottom: 4px;
}

.caption-text {
  margin: 0;
  line-height: 1.4;
}

.timestamp {
  display: block;
  font-size: 0.8em;
  color: #666;
  margin-top: 8px;
}

.comments-section {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
}

.actions-section {
  padding: 14px;
  border-top: 1px solid #dbdbdb;
  background: white;
}

.action-button {
  background: none;
  border: none;
  padding: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 6px;
}

.feather {
  width: 24px;
  height: 24px;
  stroke: currentColor;
}

.feather.liked {
  fill: #ed4956;
  stroke: #ed4956;
}

.likes-count {
  font-size: 0.9rem;
  color: #262626;
}
</style> 