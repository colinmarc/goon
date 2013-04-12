#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include "parser.h"

char *format_node(gn_ast_node_t *node) {
    char *left = "";
    char *right = "";

    if (node->num_children > 0) {
        left = format_node(node->children[0]);
    }

    if (node->num_children > 1) {
        right = format_node(node->children[1]);
    }

    char *s = malloc(1023 * sizeof(s));
    sprintf(
        s,
        "{type: %d, value: %d, left: %s, right: %s}",
        (int)node->node_type, node->value, left, right
    );

    return s;
}

int main() {
    char *s = "foo = 5\n";
    gn_ast_node_t *root = gn_parse(gn_global_context(), s, strlen(s));
    printf("%s\n", format_node(root));

    // s = "a = 1\n";
    // root = gn_parse(gn_global_context(), s, strlen(s));
    // printf("%s\n", format_node(root));
}