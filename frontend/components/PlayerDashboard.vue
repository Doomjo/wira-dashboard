<template>
  <div class="dashboard-container">
    <header class="dashboard-header">
      <h1 class="dashboard-title">WIRA Ranking Dashboard</h1>
    </header>
    
    <main class="dashboard-main">
      <!-- Search Bar -->
      <div class="search-container">
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Search players..."
          class="search-input"
          @input="handleSearch"
        />
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="loading-spinner">
        Loading...
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="error-message">
        {{ error }}
      </div>

      <!-- Players Table -->
      <div v-else class="table-container">
        <table class="players-table">
          <thead>
            <tr>
              <th>Rank</th>
              <th>Username</th>
              <th>Class</th>
              <th>Score</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(player, index) in players" :key="player.char_id">
              <td>
                <span :class="['rank-badge', index < 3 ? 'top-rank' : '']">
                  {{ calculateRank(index) }}
                </span>
              </td>
              <td>{{ player.username }}</td>
              <td>
                <span class="class-badge">
                  Class {{ player.class_id }}
                </span>
              </td>
              <td class="score">{{ player.reward_score.toLocaleString() }}</td>
            </tr>
          </tbody>
        </table>

        <!-- Pagination -->
        <div class="pagination">
          <div class="pagination-info">
            Showing {{ ((currentPage - 1) * itemsPerPage) + 1 }} to 
            {{ Math.min(currentPage * itemsPerPage, totalCount) }} of 
            {{ totalCount }} entries
          </div>
          <div class="pagination-buttons">
            <button
              @click="currentPage--"
              :disabled="currentPage === 1"
              class="pagination-button"
            >
              Previous
            </button>
            <button
              @click="currentPage++"
              :disabled="currentPage >= totalPages"
              class="pagination-button"
            >
              Next
            </button>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script>
import axios from 'axios';
import { debounce } from 'lodash';

export default {
  name: 'PlayerDashboard',
  data() {
    return {
      players: [],
      searchQuery: '',
      currentPage: 1,
      totalPages: 1,
      totalCount: 0,
      itemsPerPage: 10,
      loading: false,
      error: null
    };
  },
  watch: {
    currentPage() {
      this.fetchPlayers();
    }
  },
  methods: {
    async fetchPlayers() {
      try {
        this.loading = true;
        this.error = null;
        
        const response = await axios.get('http://localhost:8080/api/players', {
          params: {
            page: this.currentPage,
            limit: this.itemsPerPage,
            search: this.searchQuery
          }
        });

        this.players = response.data.players;
        this.totalCount = response.data.total_count;
        this.totalPages = response.data.total_pages;
      } catch (err) {
        this.error = 'Error fetching players data';
        console.error(err);
      } finally {
        this.loading = false;
      }
    },
    handleSearch: debounce(function() {
      this.currentPage = 1;
      this.fetchPlayers();
    }, 300),
    calculateRank(index) {
      return ((this.currentPage - 1) * this.itemsPerPage) + index + 1;
    }
  },
  created() {
    this.fetchPlayers();
  }
};
</script>

<style scoped>
.dashboard-container {
  min-height: 100vh;
  background-color: white;
}

.dashboard-header {
  background-color: #f0f9ff;
  padding: 1rem 0;
  text-align: center;
  border-bottom: 1px solid #e5e7eb;
}

.dashboard-title {
  font-size: 2.5rem;
  font-weight: bold;
  color: #2563eb;
}

.dashboard-main {
  max-width: 95%;
  margin: 0 auto;
  padding: 2rem 1rem;
}

.search-container {
  max-width: 32rem;
  margin: 2rem auto;
}

.search-input {
  width: 100%;
  padding: 0.75rem 1rem;
  border: 1px solid #e5e7eb;
  border-radius: 0.5rem;
  font-size: 1rem;
}

.search-input:focus {
  outline: none;
  border-color: #3b82f6;
  box-shadow: 0 0 0 2px solid #3b82f6;
}

.table-container {
  background-color: white;
  border-radius: 0.5rem;
  border: 1px solid #e5e7eb;
  overflow: hidden;
}

.players-table {
  width: 100%;
  border-collapse: collapse;
}

.players-table th,
.players-table td {
  padding: 1rem 1.5rem;
  text-align: left;
  border-bottom: 1px solid #e5e7eb;
}

.players-table th {
  background-color: #f9fafb;
  font-weight: 500;
  text-transform: uppercase;
  font-size: 0.75rem;
  color: #6b7280;
}

.players-table tr:hover {
  background-color: #f0f9ff;
}

.rank-badge {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  border-radius: 9999px;
  background-color: #f3f4f6;
  color: #4b5563;
}

.top-rank {
  background-color: #dbeafe;
  color: #1e40af;
}

.class-badge {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  border-radius: 9999px;
  background-color: #d1fae5;
  color: #065f46;
  font-size: 0.875rem;
}

.score {
  font-weight: 500;
}

.pagination {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 1.5rem;
  border-top: 1px solid #e5e7eb;
}

.pagination-info {
  font-size: 0.875rem;
  color: #4b5563;
}

.pagination-buttons {
  display: flex;
  gap: 2rem;
}

.pagination-button {
  padding: 0.5rem 1.5rem;
  border: 1px solid #e5e7eb;
  border-radius: 9999px;
  background-color: white;
  color: #4b5563;
  cursor: pointer;
}

.pagination-button:hover:not(:disabled) {
  background-color: #f9fafb;
}

.pagination-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.loading-spinner {
  text-align: center;
  padding: 3rem 0;
  color: #3b82f6;
}

.error-message {
  text-align: center;
  padding: 1rem;
  color: #ef4444;
}
</style>
