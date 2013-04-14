#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "parser.h"

gn_parser_context_t *gn_context;

static void gn_init_context() {
    gn_context = malloc(sizeof(*gn_context));

    gn_context->position = 0;
    gn_context->next_symbol = 0;
    gn_context->buffer_len = 0;
    gn_context->buffer_offset = 0;
    gn_context->error = 0;
}

static void gn_set_error(int error) {
    gn_context->error = error;
}

static gn_ast_node_t *gn_create_node(gn_ast_node_type_t node_type, int value) {
    gn_ast_node_t *node = malloc(sizeof(*node));

    node->node_type = node_type;
    node->value = value;
    node->num_children = 0;

    return node;
}

static gn_ast_node_t *gn_top() {
    return gn_context->stack[gn_context->position - 1];
}

static void gn_push(gn_ast_node_t *node)	{
    //printf("pushing %d/%d/%d\n", node->node_type, node->value, node->num_children);
    gn_context->stack[gn_context->position] = node;
    gn_context->position++;
}

static gn_ast_node_t *gn_pop() {
    if (gn_context->position <= 0) return NULL;

    gn_ast_node_t *node = gn_top();
    gn_context->position--;
    return node;
}

static void gn_reduce(gn_ast_node_type_t node_type, int num_children) {
    //printf("reducing! %d, %d\n", node_type, num_children);
    gn_ast_node_t **children = malloc(num_children * sizeof(gn_ast_node_t *));

    int i;
    for(i = num_children - 1; i >= 0; i--) {
        gn_ast_node_t *node = gn_pop();
        children[i] = node;
    }

    gn_ast_node_t *node = gn_create_node(node_type, -1);
    node->children = children;
    node->num_children = num_children;

    gn_push(node);
}

static void gn_nil_node() {
    gn_push(gn_create_node(GN_AST_NIL, -1));
}

static void gn_number_node(char *text) {
    int value = atoi(text);
    gn_push(gn_create_node(GN_AST_NUMBER, value));
}

static void gn_bool_node(int value) {
    gn_push(gn_create_node(GN_AST_BOOLEAN, value));
}

static void gn_symbol_node(char *text) {
    int i = gn_context->next_symbol++;

    int text_len = strlen(text);
    char *symbol = malloc(text_len * sizeof(*symbol));
    strncpy(symbol, text, text_len);

    gn_context->symbols[i] = symbol;

    gn_push(gn_create_node(GN_AST_SYMBOL, i));
}

static int gn_read_input(char *buffer, int *result, int max_size) {
    if(gn_context->buffer_len == 0) {
        *result = 0;
        return 0;
    }

    int num_bytes = gn_context->buffer_len;
    if (num_bytes > max_size) num_bytes = max_size;

    int i;
    for (i = 0; i < num_bytes; i++) {
        buffer[i] = gn_context->input_buffer[gn_context->buffer_offset+i];
    }

    *result = num_bytes;
    gn_context->buffer_offset += num_bytes;
    gn_context->buffer_len -= num_bytes;

    return 0;
}

static void gn_feed_input(char *new_input, int new_input_len) {
    gn_context->input_buffer = new_input;
    gn_context->buffer_len = new_input_len;
    gn_context->buffer_offset = 0;
}

#undef YY_INPUT
#define YY_INPUT(buf, result, max_size) gn_read_input(buf, &result, max_size);

#include "goon.peg.c"

// public

gn_ast_node_t *gn_parse(char *buffer, int len) {
    gn_init_context();

    gn_feed_input(buffer, len);
    yyparse();

    if (gn_context->error) {
        return NULL;
    }

    gn_ast_node_t *root = gn_pop();
    return root;
}

gn_ast_node_t *gn_child_at(gn_ast_node_t *parent, int idx) {
    return parent->children[idx];
}

char *gn_get_symbol(gn_ast_node_t *node) {
    return gn_context->symbols[node->value];
}
