import { createRouter, createWebHistory } from 'vue-router'
import HomeView from "@/views/HomeView.vue";
import GameView from "@/views/GameView.vue";
import {usePlayerStore} from "@/stores/player.store.ts";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/game',
      name: 'game',
      component: GameView,
    }
  ]
});

router.beforeEach(async (to, from) => {
  const playerStore = usePlayerStore();
  if(!playerStore.player && to.name != 'home') {
    return {
      name: 'home'
    }
  }
});

export default router
