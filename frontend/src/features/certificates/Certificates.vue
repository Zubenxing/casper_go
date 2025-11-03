<template>
  <div class="certificates-page">
    <CertificateList v-if="activeTab === 'list'" />
    <CertUpdate v-if="activeTab === 'update'" />
    <GenerateCSR v-if="activeTab === 'generate-csr'" />
    <ValidateCSR v-if="activeTab === 'validate-csr'" />
    <ValidateCert v-if="activeTab === 'validate-cert'" />
  </div>
</template>

<script setup>
import { inject } from 'vue'
import CertificateList from './CertificateList.vue'
import CertUpdate from './CertUpdate.vue'
import GenerateCSR from './GenerateCSR.vue'
import ValidateCSR from './ValidateCSR.vue'
import ValidateCert from './ValidateCert.vue'

// 从 MainLayout 注入
const activeTab = inject('certificateTab')
</script>

<style scoped>
.certificates-page {
  width: 100%;
}
</style>

