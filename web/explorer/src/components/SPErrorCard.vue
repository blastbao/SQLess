<template>
    <v-card v-if="!!error">
        <v-card-title primary-title>
            <v-flex headline text-md-center text-xs-center error--text text--darken-1>
                ERROR
            </v-flex>
        </v-card-title>
        <v-card-text>
            <v-flex text-md-center text-xs-center grey--text>
                Server said: {{ error.status }}
            </v-flex>
        </v-card-text>
    </v-card>
</template>

<script>
export default {
  props: ['error']
}
</script>
